package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"nft-auction-backend/internal/service"
)

type AuctionHandler struct {
	service *service.AuctionService
}

func NewAuctionHandler(service *service.AuctionService) *AuctionHandler {
	return &AuctionHandler{service: service}
}

// GetAuctions 获取所有拍卖
func (h *AuctionHandler) GetAuctions(c *gin.Context) {
	auctions, err := h.service.GetAllAuctions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取拍卖列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    auctions,
		"total":   len(auctions),
	})
}

// GetAuction 获取单个拍卖
func (h *AuctionHandler) GetAuction(c *gin.Context) {
	idStr := c.Param("id")
	auctionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的拍卖ID",
		})
		return
	}

	auction, err := h.service.GetAuctionByID(auctionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "拍卖不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    auction,
	})
}

// SyncAuctions 手动同步拍卖
func (h *AuctionHandler) SyncAuctions(c *gin.Context) {
	err := h.service.SyncAuctions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "同步失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "同步成功",
	})
}

// GetActiveAuctions 获取进行中的拍卖
func (h *AuctionHandler) GetActiveAuctions(c *gin.Context) {
	auctions, err := h.service.GetActiveAuctions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取进行中拍卖失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    auctions,
		"total":   len(auctions),
	})
}
