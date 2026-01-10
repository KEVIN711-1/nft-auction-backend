package service

import (
	"fmt"
	"log"
	"time"

	"nft-auction-backend/internal/contract"
	"nft-auction-backend/internal/model"

	"gorm.io/gorm"
)

type AuctionService struct {
	db       *gorm.DB
	contract *contract.ContractClient
}

func NewAuctionService(db *gorm.DB, contract *contract.ContractClient) *AuctionService {
	return &AuctionService{
		db:       db,
		contract: contract,
	}
}

// SyncAuctions 从区块链同步拍卖数据
func (s *AuctionService) SyncAuctions() error {
	log.Println("🔄 开始同步拍卖数据...")

	// 1. 获取最新区块（测试连接）
	blockNumber, err := s.contract.GetLatestBlockNumber()
	if err != nil {
		return fmt.Errorf("获取区块信息失败: %v", err)
	}
	log.Printf("当前区块高度: %d", blockNumber)

	// 2. 获取模拟数据（实际应该调用合约）
	mockAuctions, err := s.contract.GetMockAuctions()
	if err != nil {
		return fmt.Errorf("获取拍卖数据失败: %v", err)
	}

	// 3. 保存到数据库
	for _, mockAuction := range mockAuctions {
		auction := model.Auction{
			AuctionID:     mockAuction.AuctionID,
			NFTContract:   "0x742d35Cc6634C0532925a3b844Bc9e0BBd17e1f6", // 示例合约
			TokenID:       fmt.Sprintf("%d", mockAuction.AuctionID),
			Seller:        mockAuction.Seller,
			StartingPrice: mockAuction.StartingPrice.String(),
			HighestBid:    mockAuction.HighestBid.String(),
			HighestBidder: mockAuction.HighestBidder,
			StartTime:     uint64(time.Now().Add(-24 * time.Hour).Unix()), // 24小时前开始
			EndTime:       uint64(time.Now().Add(24 * time.Hour).Unix()),  // 24小时后结束
			Ended:         false,
			TxHash:        fmt.Sprintf("0x%064d", mockAuction.AuctionID),
		}

		// 检查是否已存在
		var existing model.Auction
		result := s.db.Where("auction_id = ?", mockAuction.AuctionID).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// 创建新记录
			if err := s.db.Create(&auction).Error; err != nil {
				log.Printf("创建拍卖 %d 失败: %v", mockAuction.AuctionID, err)
			} else {
				log.Printf("✅ 创建拍卖 %d", mockAuction.AuctionID)
			}
		} else {
			// 更新现有记录
			if err := s.db.Model(&existing).Updates(&auction).Error; err != nil {
				log.Printf("更新拍卖 %d 失败: %v", mockAuction.AuctionID, err)
			} else {
				log.Printf("🔄 更新拍卖 %d", mockAuction.AuctionID)
			}
		}
	}

	log.Printf("✅ 同步完成，处理了 %d 个拍卖", len(mockAuctions))
	return nil
}

// GetAllAuctions 获取所有拍卖
func (s *AuctionService) GetAllAuctions() ([]model.Auction, error) {
	var auctions []model.Auction
	result := s.db.Order("created_at DESC").Find(&auctions)
	return auctions, result.Error
}

// GetAuctionByID 根据ID获取拍卖
func (s *AuctionService) GetAuctionByID(auctionID uint64) (*model.Auction, error) {
	var auction model.Auction
	result := s.db.Where("auction_id = ?", auctionID).First(&auction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &auction, nil
}

// GetActiveAuctions 获取进行中的拍卖
func (s *AuctionService) GetActiveAuctions() ([]model.Auction, error) {
	var auctions []model.Auction
	currentTime := uint64(time.Now().Unix())
	result := s.db.Where("ended = ? AND end_time > ?", false, currentTime).
		Order("end_time ASC"). // 按结束时间升序
		Find(&auctions)
	return auctions, result.Error
}
