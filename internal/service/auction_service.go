package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"nft-auction-backend/internal/contract"
	"nft-auction-backend/internal/model"

	"gorm.io/gorm"
)

// AuctionService 处理拍卖业务逻辑
type AuctionService struct {
	db        *gorm.DB
	nftClient contract.NFTContract
}

// NewAuctionService 创建拍卖服务
func NewAuctionService(db *gorm.DB, nftClient contract.NFTContract) *AuctionService {
	return &AuctionService{
		db:        db,
		nftClient: nftClient,
	}
}

// CreateAuction 创建拍卖
func (s *AuctionService) CreateAuction(ctx context.Context, req *CreateAuctionRequest) (*model.Auction, error) {
	// 验证 NFT 所有权
	tokenID, ok := new(big.Int).SetString(req.TokenID, 10)
	if !ok {
		return nil, fmt.Errorf("invalid Token ID: %s", req.TokenID)
	}

	// 使用 NFTClient 验证所有权
	isOwner, err := s.nftClient.CheckOwner(ctx, tokenID, req.Seller)
	if err != nil {
		return nil, fmt.Errorf("failed to verify NFT ownership: %v", err)
	}
	if !isOwner {
		return nil, errors.New("seller does not own the NFT")
	}

	// 检查 NFT 是否已在其他活跃拍卖中
	var existingAuction model.Auction
	err = s.db.Where("token_id = ? AND ended = ?", req.TokenID, false).First(&existingAuction).Error
	if err == nil {
		return nil, errors.New("NFT is already in an active auction")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing auctions: %v", err)
	}

	// 验证时间
	now := time.Now()
	if req.StartTime.Before(now) {
		return nil, errors.New("start time must be in the future")
	}
	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end time must be after start time")
	}
	if req.EndTime.Sub(req.StartTime) < time.Hour {
		return nil, errors.New("auction duration must be at least 1 hour")
	}

	// 获取下一个 AuctionID（如果前端不提供的话）
	// 这里简化处理，实际情况可能需要更复杂的ID生成逻辑
	var maxID uint64
	s.db.Model(&model.Auction{}).Select("COALESCE(MAX(auction_id), 0)").Scan(&maxID)
	auctionID := maxID + 1

	// 创建拍卖 - 根据你的模型字段
	auction := &model.Auction{
		AuctionID:     auctionID,
		NFTContract:   req.NFTContract,
		TokenID:       req.TokenID,
		Seller:        strings.ToLower(req.Seller),
		StartingPrice: strconv.FormatFloat(req.StartingPrice, 'f', -1, 64),
		HighestBid:    "0", // 初始最高出价为0
		HighestBidder: "",
		StartTime:     uint64(req.StartTime.Unix()),
		EndTime:       uint64(req.EndTime.Unix()),
		Ended:         false,
		TxHash:        req.TxHash,
	}

	if err := s.db.Create(auction).Error; err != nil {
		return nil, fmt.Errorf("failed to create auction: %v", err)
	}

	return auction, nil
}

// PlaceBid 出价
func (s *AuctionService) PlaceBid(ctx context.Context, auctionID uint64, req *PlaceBidRequest) error {
	var auction model.Auction
	if err := s.db.Where("auction_id = ?", auctionID).First(&auction).Error; err != nil {
		return fmt.Errorf("auction not found: %v", err)
	}

	// 检查拍卖是否已结束
	if auction.Ended {
		return errors.New("auction has already ended")
	}

	// 检查拍卖时间
	now := time.Now()
	startTime := time.Unix(int64(auction.StartTime), 0)
	endTime := time.Unix(int64(auction.EndTime), 0)

	if now.Before(startTime) {
		return errors.New("auction has not started yet")
	}
	if now.After(endTime) {
		return errors.New("auction has ended")
	}

	// 检查出价者不是卖家
	if strings.EqualFold(req.Bidder, auction.Seller) {
		return errors.New("seller cannot bid on their own auction")
	}

	// 转换当前最高出价
	currentHighest, err := strconv.ParseFloat(auction.HighestBid, 64)
	if err != nil {
		currentHighest = 0
	}

	// 转换起拍价
	startingPrice, err := strconv.ParseFloat(auction.StartingPrice, 64)
	if err != nil {
		startingPrice = 0
	}

	// 检查出价是否高于当前最高价和起拍价
	minBid := currentHighest
	if currentHighest < startingPrice {
		minBid = startingPrice
	}
	if req.Amount <= minBid {
		return fmt.Errorf("bid must be higher than %.2f", minBid)
	}

	// 使用事务确保一致性
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 更新拍卖的最高出价
		auction.HighestBid = strconv.FormatFloat(req.Amount, 'f', -1, 64)
		auction.HighestBidder = req.Bidder
		if err := tx.Save(&auction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to place bid: %v", err)
	}

	return nil
}

// EndAuction 结束拍卖
func (s *AuctionService) EndAuction(ctx context.Context, auctionID uint64, callerAddress string) error {
	var auction model.Auction
	if err := s.db.Where("auction_id = ?", auctionID).First(&auction).Error; err != nil {
		return fmt.Errorf("auction not found: %v", err)
	}

	// 检查权限：只有卖家或拍卖已结束才能结束拍卖
	now := time.Now()
	endTime := time.Unix(int64(auction.EndTime), 0)
	isSeller := strings.EqualFold(callerAddress, auction.Seller)
	isEnded := now.After(endTime)

	if !isSeller && !isEnded {
		return errors.New("only seller can end auction before end time")
	}

	if auction.Ended {
		return errors.New("auction has already ended")
	}

	// 更新拍卖状态
	auction.Ended = true
	if err := s.db.Save(&auction).Error; err != nil {
		return fmt.Errorf("failed to end auction: %v", err)
	}

	// 如果有最高出价者，记录转移信息
	if auction.HighestBidder != "" {
		log.Printf("Auction %d ended. NFT %s should be transferred from %s to %s",
			auction.AuctionID, auction.TokenID, auction.Seller, auction.HighestBidder)

		// TODO: 实际调用智能合约转移 NFT
		// tokenID, _ := new(big.Int).SetString(auction.TokenID, 10)
		// from := common.HexToAddress(auction.Seller)
		// to := common.HexToAddress(auction.HighestBidder)
		// err := s.nftClient.TransferFrom(ctx, from, to, tokenID)
	}

	return nil
}

// GetAuction 获取拍卖详情
func (s *AuctionService) GetAuction(ctx context.Context, auctionID uint64) (*model.Auction, error) {
	var auction model.Auction
	if err := s.db.Where("auction_id = ?", auctionID).First(&auction).Error; err != nil {
		return nil, fmt.Errorf("auction not found: %v", err)
	}
	return &auction, nil
}

// ListAuctions 获取拍卖列表
func (s *AuctionService) ListAuctions(ctx context.Context, ended *bool, limit, offset int) ([]model.Auction, error) {
	var auctions []model.Auction
	query := s.db

	if ended != nil {
		query = query.Where("ended = ?", *ended)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Order("created_at DESC").Find(&auctions).Error; err != nil {
		return nil, fmt.Errorf("failed to list auctions: %v", err)
	}

	return auctions, nil
}

// GetActiveAuctions 获取活跃拍卖（未结束）
func (s *AuctionService) GetActiveAuctions(ctx context.Context, limit, offset int) ([]model.Auction, error) {
	ended := false
	return s.ListAuctions(ctx, &ended, limit, offset)
}

// GetEndedAuctions 获取已结束拍卖
func (s *AuctionService) GetEndedAuctions(ctx context.Context, limit, offset int) ([]model.Auction, error) {
	ended := true
	return s.ListAuctions(ctx, &ended, limit, offset)
}

// SyncAuctions 同步拍卖数据（如果还需要的话）
func (s *AuctionService) SyncAuctions() error {
	// 这里可以实现从链上同步拍卖数据的逻辑
	// 目前可以先留空或实现基础逻辑
	log.Println("SyncAuctions called - placeholder implementation")
	return nil
}

// Request types
type CreateAuctionRequest struct {
	NFTContract   string    `json:"nft_contract"`
	TokenID       string    `json:"token_id"`
	Seller        string    `json:"seller"`
	StartingPrice float64   `json:"starting_price"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	TxHash        string    `json:"tx_hash,omitempty"`
}

type PlaceBidRequest struct {
	Bidder string  `json:"bidder"`
	Amount float64 `json:"amount"`
}
