// auction_service.go
package service

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"gorm.io/gorm"

	"nft-auction-backend/internal/contract"
	"nft-auction-backend/internal/model"
)

// AuctionService 拍卖服务（只读，不包含需要gas的操作）
type AuctionService struct {
	DB              *gorm.DB
	AuctionContract contract.AuctionContract
}

// NewAuctionService 创建拍卖服务
func NewAuctionService(db *gorm.DB, auctionContract contract.AuctionContract) *AuctionService {
	return &AuctionService{
		DB:              db,
		AuctionContract: auctionContract,
	}
}

// ==================== 数据库操作 ====================

// SaveAuction 保存或更新拍卖到数据库
func (s *AuctionService) SaveAuction(auction *model.Auction) error {
	if auction == nil {
		return fmt.Errorf("拍卖信息为空")
	}

	var existing model.Auction
	result := s.DB.Where("auction_id = ?", auction.AuctionID).First(&existing)
	now := time.Now()

	if result.Error != nil {
		// 新记录
		auction.CreatedAt = now
		auction.UpdatedAt = now

		if err := s.DB.Create(auction).Error; err != nil {
			return fmt.Errorf("创建拍卖失败: %v", err)
		}
		log.Printf("✅ 新增拍卖 #%d", auction.AuctionID)
	} else {
		// 更新现有记录
		existing.NFTContract = auction.NFTContract
		existing.TokenID = auction.TokenID
		existing.Seller = auction.Seller
		existing.StartingPrice = auction.StartingPrice
		existing.HighestBid = auction.HighestBid
		existing.HighestBidder = auction.HighestBidder
		existing.StartTime = auction.StartTime
		existing.EndTime = auction.EndTime
		existing.Ended = auction.Ended
		existing.Status = auction.Status
		existing.UpdatedAt = now

		if err := s.DB.Save(&existing).Error; err != nil {
			return fmt.Errorf("更新拍卖失败: %v", err)
		}
		log.Printf("🔄 更新拍卖 #%d", auction.AuctionID)
	}

	return nil
}

// SaveBidHistory 保存出价历史记录
func (s *AuctionService) SaveBidHistory(bid *model.BidHistory) error {
	if bid == nil {
		return fmt.Errorf("出价记录为空")
	}

	// 检查是否已存在（根据交易哈希）
	var existing model.BidHistory
	if err := s.DB.Where("tx_hash = ?", bid.TxHash).First(&existing).Error; err == nil {
		// 已存在，更新
		existing.Amount = bid.Amount
		existing.Status = bid.Status
		existing.BlockNumber = bid.BlockNumber
		existing.BlockTime = bid.BlockTime
		existing.UpdatedAt = time.Now()

		if err := s.DB.Save(&existing).Error; err != nil {
			return fmt.Errorf("更新出价记录失败: %v", err)
		}
		return nil
	}

	// 新记录
	now := time.Now()
	bid.CreatedAt = now
	bid.UpdatedAt = now

	if err := s.DB.Create(bid).Error; err != nil {
		return fmt.Errorf("创建出价记录失败: %v", err)
	}

	log.Printf("✅ 保存出价记录: AuctionID=%d, Bidder=%s", bid.AuctionID, bid.Bidder)
	return nil
}

// ==================== 链上查询方法 ====================

// GetAuctionFromChain 从区块链获取拍卖信息（适配你的接口）
func (s *AuctionService) GetAuctionFromChain(ctx context.Context, auctionID uint64) (*model.Auction, error) {
	// 使用 GetAuctionInfo 方法获取拍卖信息
	seller, duration, startPrice, startTime, ended, highestBidder, highestBid,
		nftContract, tokenId, _, _, _, err :=
		s.AuctionContract.GetAuctionInfo(ctx, big.NewInt(int64(auctionID)))

	if err != nil {
		return nil, fmt.Errorf("从链上获取拍卖失败: %v", err)
	}

	// 计算结束时间
	endTime := big.NewInt(0)
	if startTime != nil && duration != nil {
		endTime = new(big.Int).Add(startTime, duration)
	}

	// 判断状态
	status := "active"
	if ended {
		status = "ended"
	} else if time.Now().Unix() > endTime.Int64() {
		status = "expired"
	} else if auctionID == 0 { // 特殊处理拍卖ID为0的情况
		status = "active"
	}

	auction := &model.Auction{
		AuctionID:     auctionID,
		NFTContract:   nftContract.Hex(),
		TokenID:       tokenId.String(),
		Seller:        seller.Hex(),
		StartingPrice: startPrice.String(),
		HighestBid:    highestBid.String(),
		HighestBidder: highestBidder.Hex(),
		StartTime:     uint64(startTime.Int64()),
		EndTime:       uint64(endTime.Int64()),
		Ended:         ended,
		Status:        status,
	}

	return auction, nil
}

// SyncAllAuctions 同步所有拍卖数据到数据库
func (s *AuctionService) SyncAllAuctions(ctx context.Context) error {
	// 获取拍卖总数
	count, err := s.AuctionContract.GetAuctionCount(ctx)
	if err != nil {
		return fmt.Errorf("获取拍卖数量失败: %v", err)
	}

	log.Printf("开始同步拍卖数据，链上拍卖总数: %d", count.Int64())

	successCount := 0
	// 从0开始，因为你的拍卖ID从0开始
	for i := int64(0); i < count.Int64(); i++ {
		auctionID := uint64(i)

		// 从链上获取拍卖信息
		auction, err := s.GetAuctionFromChain(ctx, auctionID)
		if err != nil {
			log.Printf("❌ 获取拍卖 #%d 信息失败: %v", auctionID, err)
			continue
		}

		// 保存到数据库
		if err := s.SaveAuction(auction); err == nil {
			successCount++
			log.Printf("✅ 同步拍卖 #%d: NFT=%s/%s, 最高出价=%s",
				auctionID, auction.NFTContract, auction.TokenID, auction.HighestBid)
		} else {
			log.Printf("❌ 保存拍卖 #%d 失败: %v", auctionID, err)
		}
	}

	log.Printf("✅ 拍卖同步完成，成功同步: %d/%d", successCount, count.Int64())
	return nil
}

// ==================== 查询方法 ====================

// GetAuctionByID 根据数据库ID获取拍卖
func (s *AuctionService) GetAuctionByID(id uint) (*model.Auction, error) {
	var auction model.Auction
	result := s.DB.First(&auction, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &auction, nil
}

// GetAuctionByAuctionID 根据链上AuctionID获取拍卖
func (s *AuctionService) GetAuctionByAuctionID(auctionID uint64) (*model.Auction, error) {
	var auction model.Auction
	result := s.DB.Where("auction_id = ?", auctionID).First(&auction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &auction, nil
}

// GetAuctionByTxHash 根据交易哈希获取拍卖（用于前端提交后查询）
func (s *AuctionService) GetAuctionByTxHash(txHash string) (*model.Auction, error) {
	var auction model.Auction
	result := s.DB.Where("tx_hash = ?", txHash).First(&auction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &auction, nil
}

// GetActiveAuctions 获取所有活跃拍卖
func (s *AuctionService) GetActiveAuctions() ([]model.Auction, error) {
	var auctions []model.Auction
	currentTime := uint64(time.Now().Unix())
	log.Printf("✅ ----currentTime=%d ", currentTime)

	result := s.DB.Where("ended = ? AND end_time > ?", false, currentTime).
		Order("created_at DESC").
		Find(&auctions)

	if result.Error != nil {
		return nil, result.Error
	}
	return auctions, nil
}

// GetAuctionBids 获取拍卖的出价历史
func (s *AuctionService) GetAuctionBids(auctionID uint64, page, pageSize int) ([]model.BidHistory, int64, error) {
	var bids []model.BidHistory
	var total int64

	query := s.DB.Model(&model.BidHistory{}).Where("auction_id = ?", auctionID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&bids).Error

	if err != nil {
		return nil, 0, err
	}

	return bids, total, nil
}

// UpdateAuctionFromChain 从链上更新单个拍卖信息（事件监听器调用）
func (s *AuctionService) UpdateAuctionFromChain(auctionID uint64) error {
	ctx := context.Background()

	auction, err := s.GetAuctionFromChain(ctx, auctionID)
	if err != nil {
		return fmt.Errorf("获取链上拍卖信息失败: %v", err)
	}

	return s.SaveAuction(auction)
}

// ValidateAuctionExists 验证拍卖是否存在（只读检查）
func (s *AuctionService) ValidateAuctionExists(ctx context.Context, auctionID uint64) (bool, error) {
	_, _, _, _, _, _, _, _, _, _, _, _, err :=
		s.AuctionContract.GetAuctionInfo(ctx, big.NewInt(int64(auctionID)))

	if err != nil {
		// 检查是否是"拍卖不存在"的错误
		if err.Error() == "execution reverted" ||
			err.Error() == "auction does not exist" ||
			err.Error() == "Not exist" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetAuctionCount 获取拍卖总数
func (s *AuctionService) GetAuctionCount(ctx context.Context) (int64, error) {
	count, err := s.AuctionContract.GetAuctionCount(ctx)
	if err != nil {
		return 0, err
	}
	return count.Int64(), nil
}

// GetContractInfo 获取合约信息
func (s *AuctionService) GetContractInfo(ctx context.Context) (map[string]interface{}, error) {
	info := make(map[string]interface{})

	// 获取拍卖总数
	count, err := s.AuctionContract.GetAuctionCount(ctx)
	if err == nil {
		info["auction_count"] = count.Int64()
	}

	// 获取合约地址
	info["contract_address"] = s.AuctionContract.GetContractAddress().Hex()

	// 获取一些活跃拍卖作为示例
	activeAuctions, _ := s.GetActiveAuctions()
	info["active_auctions"] = len(activeAuctions)

	return info, nil
}
