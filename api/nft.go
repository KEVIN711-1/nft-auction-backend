package api

import (
	"net/http"
	"strings"
	"time"

	"nft-auction-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type NFTHandler struct {
	service *service.NFTService
}

func NewNFTHandler(service *service.NFTService) *NFTHandler {
	return &NFTHandler{service: service}
}

// GetNFTInfo 获取NFT信息
func (h *NFTHandler) GetNFTInfo(c *gin.Context) {

	tokenID := c.Param("id")
	if tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token ID is required",
		})
		return
	}

	// 调用服务层获取NFT信息
	contractAddr := h.service.GetContractAddress()

	info, err := h.service.GetNFT(contractAddr, tokenID)
	if err != nil {
		// 根据错误类型返回不同的状态码
		if strings.Contains(err.Error(), "does not exist") ||
			strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "NFT not found",
				"details": err.Error(),
			})
		} else if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid token ID",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get NFT info",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetNFTOwner 获取NFT所有者
func (h *NFTHandler) GetNFTOwner(c *gin.Context) {
	ctx := c.Request.Context()

	tokenID := c.Param("id")
	if tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token ID is required",
		})
		return
	}

	owner, err := h.service.GetOwner(ctx, tokenID)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") ||
			strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "NFT not found",
				"details": err.Error(),
			})
		} else if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid token ID",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get NFT owner",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token_id": tokenID,
		"owner":    owner,
	})
}

// ValidateOwnership 验证NFT所有权
func (h *NFTHandler) ValidateOwnership(c *gin.Context) {
	ctx := c.Request.Context()

	tokenID := c.Param("id")
	address := c.Param("address")

	if tokenID == "" || address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token ID and address are required",
		})
		return
	}

	// 清理地址格式
	address = strings.ToLower(strings.TrimSpace(address))

	isOwner, err := h.service.ValidateOwnership(ctx, tokenID, address)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") ||
			strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "NFT not found",
				"details": err.Error(),
			})
		} else if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid parameters",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to validate ownership",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token_id": tokenID,
		"address":  address,
		"is_owner": isOwner,
	})
}

// SyncAuctions 手动同步拍卖数据（管理用）
func (h *NFTHandler) SyncNFTInfo(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.service.SyncAllNFTs(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "同步失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "拍卖数据同步完成",
		"timestamp": time.Now().Unix(),
	})
}
