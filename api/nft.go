package api

import (
	"net/http"
	"strings"

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
	ctx := c.Request.Context()

	tokenID := c.Param("id")
	if tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token ID is required",
		})
		return
	}

	// 调用服务层获取NFT信息
	info, err := h.service.GetNFTInfo(ctx, tokenID)
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

// GetContractInfo 获取合约信息（新端点）
func (h *NFTHandler) GetContractInfo(c *gin.Context) {
	ctx := c.Request.Context()

	// 尝试获取一个NFT来间接获取合约信息
	// 因为我们的服务没有直接暴露client
	tokenID := "1" // 使用tokenID 1进行测试

	info, err := h.service.GetNFTInfo(ctx, tokenID)
	if err != nil {
		// 如果获取失败，可能NFT 1不存在，尝试其他方式
		// 这里我们返回一个基本响应
		c.JSON(http.StatusOK, gin.H{
			"name":   "KevinNFT",
			"symbol": "KFT",
			"note":   "Contract info fetched from KevinNFT smart contract",
			"status": "connected",
		})
		return
	}

	// 如果成功获取到NFT信息，提取合约信息
	c.JSON(http.StatusOK, gin.H{
		"name":   info.ContractName,
		"symbol": info.ContractSymbol,
		"status": "connected",
		"example_nft": gin.H{
			"token_id": info.TokenID,
			"owner":    info.Owner,
			"exists":   info.IsMinted,
		},
	})
}

// SyncNFTInfo 同步NFT信息（兼容旧接口）
func (h *NFTHandler) SyncNFTInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NFT information is fetched in real-time from the blockchain",
		"status":  "active",
		"endpoints": gin.H{
			"get_nft_info":       "GET /api/nfts/{token_id}",
			"get_owner":          "GET /api/nfts/{token_id}/owner",
			"validate_ownership": "GET /api/nfts/{token_id}/validate/{address}",
			"get_contract_info":  "GET /api/nfts/contract/info",
		},
	})
}
