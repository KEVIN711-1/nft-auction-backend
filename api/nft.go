// api/nft.go
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nft-auction-backend/internal/service"
)

type NFTHandler struct {
	service *service.NFTService
}

func NewNFTHandler(service *service.NFTService) *NFTHandler {
	return &NFTHandler{service: service}
}

// GetNFTInfo 获取NFT信息
func (h *NFTHandler) GetNFTInfo(c *gin.Context) {
	nftInfo, err := h.service.GetNFTInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取NFT信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    nftInfo,
	})
}

// SyncNFTInfo 手动同步NFT信息
func (h *NFTHandler) SyncNFTInfo(c *gin.Context) {
	err := h.service.SyncNFTInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "同步NFT信息失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "同步成功",
	})
}
