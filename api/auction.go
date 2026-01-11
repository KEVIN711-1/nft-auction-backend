package api

import (
	"net/http"
	"strconv"
	"strings"

	"nft-auction-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuctionHandler struct {
	service *service.AuctionService
}

func NewAuctionHandler(service *service.AuctionService) *AuctionHandler {
	return &AuctionHandler{service: service}
}

// GetAuctions 获取拍卖列表
func (h *AuctionHandler) GetAuctions(c *gin.Context) {
	ctx := c.Request.Context()

	// 获取查询参数
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var ended *bool
	if endedStr := c.Query("ended"); endedStr != "" {
		endedVal, err := strconv.ParseBool(endedStr)
		if err == nil {
			ended = &endedVal
		}
	}

	auctions, err := h.service.ListAuctions(ctx, ended, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch auctions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auctions": auctions,
		"count":    len(auctions),
		"limit":    limit,
		"offset":   offset,
	})
}

// GetActiveAuctions 获取活跃拍卖
func (h *AuctionHandler) GetActiveAuctions(c *gin.Context) {
	ctx := c.Request.Context()

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	auctions, err := h.service.GetActiveAuctions(ctx, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch active auctions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auctions": auctions,
		"count":    len(auctions),
	})
}

// GetAuction 获取单个拍卖详情
func (h *AuctionHandler) GetAuction(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	auctionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid auction ID",
		})
		return
	}

	auction, err := h.service.GetAuction(ctx, auctionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Auction not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, auction)
}

// CreateAuction 创建拍卖
func (h *AuctionHandler) CreateAuction(c *gin.Context) {
	ctx := c.Request.Context()

	var req service.CreateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// 验证必填字段
	if req.NFTContract == "" || req.TokenID == "" || req.Seller == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required fields: nft_contract, token_id, seller",
		})
		return
	}

	auction, err := h.service.CreateAuction(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create auction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, auction)
}

// PlaceBid 出价
func (h *AuctionHandler) PlaceBid(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	auctionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid auction ID",
		})
		return
	}

	var req service.PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Bidder == "" || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required fields: bidder, amount (must be positive)",
		})
		return
	}

	err = h.service.PlaceBid(ctx, auctionID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "auction not found" {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "auction has already ended") ||
			strings.Contains(err.Error(), "auction has not started yet") ||
			strings.Contains(err.Error(), "seller cannot bid") ||
			strings.Contains(err.Error(), "bid must be higher") {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{
			"error":   "Failed to place bid",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Bid placed successfully",
		"auction_id": auctionID,
		"bidder":     req.Bidder,
		"amount":     req.Amount,
	})
}

// EndAuction 结束拍卖
func (h *AuctionHandler) EndAuction(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	auctionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid auction ID",
		})
		return
	}

	// 获取调用者地址（从请求头或认证信息中）
	callerAddress := c.GetHeader("X-User-Address")
	if callerAddress == "" {
		// 如果没有从头部获取，可以从请求体中获取
		var req struct {
			CallerAddress string `json:"caller_address"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing caller address",
			})
			return
		}
		callerAddress = req.CallerAddress
	}

	if callerAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Caller address is required",
		})
		return
	}

	err = h.service.EndAuction(ctx, auctionID, callerAddress)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "auction not found" {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "only seller can end auction") ||
			strings.Contains(err.Error(), "auction has already ended") {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{
			"error":   "Failed to end auction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Auction ended successfully",
		"auction_id": auctionID,
	})
}

// SyncAuctions 同步拍卖数据（兼容旧接口）
func (h *AuctionHandler) SyncAuctions(c *gin.Context) {
	err := h.service.SyncAuctions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to sync auctions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sync completed",
	})
}
