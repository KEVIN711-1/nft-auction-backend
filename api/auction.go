package api

import (
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nft-auction-backend/internal/model"
	"nft-auction-backend/internal/service"
)

type AuctionHandler struct {
	service *service.AuctionService
}

func NewAuctionHandler(auctionService *service.AuctionService) *AuctionHandler {
	return &AuctionHandler{
		service: auctionService,
	}
}

// GetAuctions 获取所有拍卖（分页）
func (h *AuctionHandler) GetAuctions(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 直接从数据库查询（因为服务层没有 ListAuctions 方法）
	var auctions []model.Auction
	var total int64

	// 获取总数
	h.service.DB.Model(&model.Auction{}).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	h.service.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&auctions)

	c.JSON(http.StatusOK, gin.H{
		"data": auctions,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetActiveAuctions 获取进行中的拍卖
func (h *AuctionHandler) GetActiveAuctions(c *gin.Context) {
	auctions, err := h.service.GetActiveAuctions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auctions": auctions,
		"count":    len(auctions),
	})
}

// GetAuction 获取单个拍卖
func (h *AuctionHandler) GetAuction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的拍卖ID"})
		return
	}

	// 使用 GetAuctionByID（注意：参数是 uint，不是 uint64）
	auction, err := h.service.GetAuctionByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "拍卖不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, auction)
}

// CreateAuctionRequest 创建拍卖请求结构体
type CreateAuctionRequest struct {
	TokenID       string `json:"token_id" binding:"required"`
	StartingPrice string `json:"starting_price" binding:"required"`
	Duration      uint64 `json:"duration" binding:"required"`
	Seller        string `json:"seller" binding:"required"`
	NFTContract   string `json:"nft_contract" binding:"required"`
}

// CreateAuction 创建拍卖
func (h *AuctionHandler) CreateAuction(c *gin.Context) {
	var req CreateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析起拍价
	startingPrice, ok := new(big.Int).SetString(req.StartingPrice, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的价格格式"})
		return
	}

	// 生成AuctionID（这里简化处理，实际应该从链上获取或使用自增）
	// 先查询当前最大AuctionID
	var maxAuctionID uint64
	h.service.DB.Model(&model.Auction{}).Select("COALESCE(MAX(auction_id), 0)").Scan(&maxAuctionID)
	auctionID := maxAuctionID + 1

	auction := &model.Auction{
		AuctionID:     auctionID,
		NFTContract:   req.NFTContract,
		TokenID:       req.TokenID,
		Seller:        req.Seller,
		StartingPrice: startingPrice.String(),
		StartTime:     uint64(time.Now().Unix()),
		EndTime:       uint64(time.Now().Unix()) + req.Duration,
		Ended:         false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 调用服务层创建拍卖
	if err := h.service.CreateAuction(auction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "拍卖创建成功",
		"auction": auction,
	})
}

// PlaceBidRequest 出价请求结构体
type PlaceBidRequest struct {
	Bidder string `json:"bidder" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

// PlaceBid 出价
func (h *AuctionHandler) PlaceBid(c *gin.Context) {
	idStr := c.Param("id")
	auctionID, err := strconv.ParseUint(idStr, 10, 32) // 转换为uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的拍卖ID"})
		return
	}

	var req PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析出价金额
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的金额格式"})
		return
	}

	// 注意：服务层的PlaceBid参数是(uint, string, *big.Int)
	if err := h.service.PlaceBid(uint(auctionID), req.Bidder, amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "出价成功",
		"auction_id": auctionID,
		"bidder":     req.Bidder,
		"amount":     amount.String(),
	})
}

// EndAuction 结束拍卖
func (h *AuctionHandler) EndAuction(c *gin.Context) {
	idStr := c.Param("id")
	auctionID, err := strconv.ParseUint(idStr, 10, 32) // 转换为uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的拍卖ID"})
		return
	}

	if err := h.service.EndAuction(uint(auctionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "拍卖结束成功",
		"auction_id": auctionID,
	})
}

// SyncAuctions 同步拍卖数据
func (h *AuctionHandler) SyncAuctions(c *gin.Context) {
	if err := h.service.SyncAuctions(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "拍卖数据同步完成",
		"timestamp": time.Now().Unix(),
	})
}
