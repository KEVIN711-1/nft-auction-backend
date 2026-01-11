package service

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"

	"nft-auction-backend/internal/contract"
	"nft-auction-backend/internal/model"
)

// AuctionService 拍卖服务
type AuctionService struct {
	DB              *gorm.DB
	NFTContract     contract.NFTContract
	AuctionContract contract.AuctionContract
}

// NewAuctionService 创建拍卖服务
func NewAuctionService(db *gorm.DB, nftContract contract.NFTContract, auctionContract contract.AuctionContract) *AuctionService {
	return &AuctionService{
		DB:              db,
		NFTContract:     nftContract,
		AuctionContract: auctionContract,
	}
}

// CreateAuction 创建拍卖（包含链上和数据库）
func (s *AuctionService) CreateAuction(auction *model.Auction) error {
	ctx := context.Background()

	// 1. 验证NFT所有权
	if s.NFTContract != nil {
		// 解析TokenID（string -> int64）
		tokenIDInt, err := strconv.ParseInt(auction.TokenID, 10, 64)
		if err != nil {
			return fmt.Errorf("无效的TokenID格式: %s, 错误: %v", auction.TokenID, err)
		}

		tokenID := big.NewInt(tokenIDInt)
		isOwner, err := s.NFTContract.CheckOwner(ctx, tokenID, auction.Seller)
		if err != nil {
			return fmt.Errorf("验证NFT所有权失败: %v", err)
		}
		if !isOwner {
			return fmt.Errorf("用户 %s 不是 NFT #%s 的所有者", auction.Seller, auction.TokenID)
		}
	}

	// 2. 调用拍卖合约创建拍卖（链上）
	if s.AuctionContract != nil && s.AuctionContract.IsActive() {
		// 将字符串参数转换为合约需要的类型
		duration := big.NewInt(int64(auction.EndTime - auction.StartTime)) // 计算持续时间

		startPrice, ok := new(big.Int).SetString(auction.StartingPrice, 10)
		if !ok {
			return fmt.Errorf("无效的起拍价格式: %s", auction.StartingPrice)
		}

		tokenIDInt, _ := strconv.ParseInt(auction.TokenID, 10, 64)
		nftAddress := common.HexToAddress(auction.NFTContract)

		// 调用合约创建拍卖（ETH拍卖）
		err := s.AuctionContract.CreateAuctionETH(ctx, duration, startPrice, nftAddress, big.NewInt(tokenIDInt))
		if err != nil {
			log.Printf("警告: 链上创建拍卖失败（继续保存到数据库）: %v", err)
			// 不返回错误，继续保存到数据库（模拟模式或测试时）
		}
	}

	// 3. 保存到数据库
	auction.CreatedAt = time.Now()
	auction.UpdatedAt = time.Now()

	// 如果Ended字段未设置，默认为false
	// 如果EndTime未设置，根据StartTime和默认持续时间计算
	if auction.EndTime == 0 && auction.StartTime > 0 {
		auction.EndTime = auction.StartTime + 86400 // 默认24小时
	}

	result := s.DB.Create(auction)
	if result.Error != nil {
		return result.Error
	}

	log.Printf("拍卖创建成功: 数据库ID=%d, AuctionID=%d, TokenID=%s",
		auction.ID, auction.AuctionID, auction.TokenID)
	return nil
}

// SyncAuctions 同步链上拍卖数据到数据库
func (s *AuctionService) SyncAuctions() error {
	if s.AuctionContract == nil || !s.AuctionContract.IsActive() {
		log.Println("拍卖合约未激活，跳过同步")
		return nil
	}

	ctx := context.Background()

	// 获取拍卖总数
	count, err := s.AuctionContract.GetAuctionCount(ctx)
	if err != nil {
		return fmt.Errorf("获取拍卖数量失败: %v", err)
	}

	log.Printf("开始同步拍卖数据，链上拍卖总数: %s", count.String())

	// 遍历所有拍卖
	for i := uint64(0); i < count.Uint64(); i++ {
		auctionID := big.NewInt(int64(i))

		// 获取拍卖信息
		// seller, duration, startPrice, startTime, ended, highestBidder, highestBid,
		// 	nftContract, tokenId, tokenAddress, bidTokenAmount, timeRemaining, err :=
		// 	s.AuctionContract.GetAuctionInfo(ctx, auctionID)
		// 获取拍卖信息 - 使用下划线忽略不需要的返回值
		seller, duration, startPrice, startTime, ended, highestBidder, highestBid,
			nftContract, tokenId, _, _, _, err :=
			s.AuctionContract.GetAuctionInfo(ctx, auctionID)

		if err != nil {
			log.Printf("获取拍卖 #%d 信息失败: %v", i, err)
			continue
		}

		// 构建拍卖记录 - 根据你的模型字段
		auction := &model.Auction{
			AuctionID:     i,
			NFTContract:   nftContract.Hex(),
			TokenID:       tokenId.String(),
			Seller:        seller.Hex(),
			StartingPrice: startPrice.String(),
			HighestBid:    highestBid.String(),
			HighestBidder: highestBidder.Hex(),
			StartTime:     startTime.Uint64(),
			Ended:         ended,
		}

		// 计算结束时间
		if startTime != nil && duration != nil {
			auction.EndTime = startTime.Uint64() + duration.Uint64()
		}

		// 检查并保存到数据库
		s.saveOrUpdateAuction(auction, i)
	}

	log.Printf("拍卖数据同步完成，处理了 %s 个拍卖", count.String())
	return nil
}

// saveOrUpdateAuction 保存或更新拍卖到数据库
func (s *AuctionService) saveOrUpdateAuction(auction *model.Auction, auctionID uint64) {
	var existingAuction model.Auction
	result := s.DB.Where("auction_id = ?", auctionID).First(&existingAuction)

	now := time.Now()

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新记录
		auction.CreatedAt = now
		auction.UpdatedAt = now
		if err := s.DB.Create(auction).Error; err != nil {
			log.Printf("创建拍卖 #%d 失败: %v", auctionID, err)
		} else {
			log.Printf("新增拍卖 #%d", auctionID)
		}
	} else {
		// 更新现有记录
		existingAuction.NFTContract = auction.NFTContract
		existingAuction.TokenID = auction.TokenID
		existingAuction.Seller = auction.Seller
		existingAuction.StartingPrice = auction.StartingPrice
		existingAuction.HighestBid = auction.HighestBid
		existingAuction.HighestBidder = auction.HighestBidder
		existingAuction.StartTime = auction.StartTime
		existingAuction.EndTime = auction.EndTime
		existingAuction.Ended = auction.Ended
		existingAuction.UpdatedAt = now

		if err := s.DB.Save(&existingAuction).Error; err != nil {
			log.Printf("更新拍卖 #%d 失败: %v", auctionID, err)
		} else {
			log.Printf("更新拍卖 #%d", auctionID)
		}
	}
}

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

// GetActiveAuctions 获取所有活跃拍卖（未结束且未过期）
func (s *AuctionService) GetActiveAuctions() ([]model.Auction, error) {
	var auctions []model.Auction
	currentTime := uint64(time.Now().Unix())

	result := s.DB.Where("ended = ? AND end_time > ?", false, currentTime).Find(&auctions)
	if result.Error != nil {
		return nil, result.Error
	}
	return auctions, nil
}

// PlaceBid 出价（链上+数据库）
func (s *AuctionService) PlaceBid(auctionID uint, bidder string, amount *big.Int) error {
	// 1. 获取拍卖信息
	var auction model.Auction
	if err := s.DB.First(&auction, auctionID).Error; err != nil {
		return err
	}

	// 2. 调用链上合约出价
	ctx := context.Background()

	// 根据你的合约逻辑，这里需要确定是ETH还是ERC20拍卖
	// 暂时假设是ETH拍卖
	err := s.AuctionContract.PlaceBidETH(ctx, big.NewInt(int64(auction.AuctionID)), amount)
	if err != nil {
		return fmt.Errorf("链上出价失败: %v", err)
	}

	// 3. 更新数据库
	auction.HighestBid = amount.String()
	auction.HighestBidder = bidder
	auction.UpdatedAt = time.Now()

	return s.DB.Save(&auction).Error
}

// EndAuction 结束拍卖
func (s *AuctionService) EndAuction(auctionID uint) error {
	// 1. 获取拍卖信息
	var auction model.Auction
	if err := s.DB.First(&auction, auctionID).Error; err != nil {
		return err
	}

	// 2. 调用链上合约结束拍卖
	ctx := context.Background()
	err := s.AuctionContract.EndAuction(ctx, big.NewInt(int64(auction.AuctionID)))
	if err != nil {
		return fmt.Errorf("链上结束拍卖失败: %v", err)
	}

	// 3. 更新数据库状态
	auction.Ended = true
	auction.UpdatedAt = time.Now()

	return s.DB.Save(&auction).Error
}

// UpdateAuctionFromChain 从链上更新单个拍卖信息
func (s *AuctionService) UpdateAuctionFromChain(auctionID uint64) error {
	ctx := context.Background()

	// 从链上获取最新信息
	seller, duration, startPrice, startTime, ended, highestBidder, highestBid,
		nftContract, tokenId, _, _, _, err :=
		s.AuctionContract.GetAuctionInfo(ctx, big.NewInt(int64(auctionID)))

	if err != nil {
		return err
	}

	// 更新数据库
	var auction model.Auction
	result := s.DB.Where("auction_id = ?", auctionID).First(&auction)
	if result.Error != nil {
		return result.Error
	}

	auction.NFTContract = nftContract.Hex()
	auction.TokenID = tokenId.String()
	auction.Seller = seller.Hex()
	auction.StartingPrice = startPrice.String()
	auction.HighestBid = highestBid.String()
	auction.HighestBidder = highestBidder.Hex()
	auction.StartTime = startTime.Uint64()
	auction.Ended = ended
	auction.UpdatedAt = time.Now()

	// 计算结束时间
	if startTime != nil && duration != nil {
		auction.EndTime = startTime.Uint64() + duration.Uint64()
	}

	return s.DB.Save(&auction).Error
}
