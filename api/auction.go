// api/auction_handler.go
package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

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

// ==================== 查询API（保留）====================

// GetAuctions 获取所有拍卖（分页，支持过滤）
func (h *AuctionHandler) GetAuctions(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	status := c.Query("status") // 按状态过滤：active, ended, all
	seller := c.Query("seller") // 按卖家过滤

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 构建查询
	var auctions []model.Auction
	var total int64

	query := h.service.DB.Model(&model.Auction{})

	// 状态过滤
	if status != "" && status != "all" {
		if status == "active" {
			query = query.Where("ended = ? AND status = ?", false, "active")
		} else if status == "ended" {
			query = query.Where("ended = ?", true)
		} else {
			query = query.Where("status = ?", status)
		}
	}

	// 卖家过滤
	if seller != "" {
		query = query.Where("seller = ?", seller)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&auctions)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"auctions": auctions,
			"pagination": gin.H{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetActiveAuctions 获取进行中的拍卖
func (h *AuctionHandler) GetActiveAuctions(c *gin.Context) {
	ctx := c.Request.Context()

	auctions, err := h.service.GetActiveAuctions(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取活跃拍卖失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"auctions": auctions,
			"count":    len(auctions),
		},
	})
}

// GetAuction 获取单个拍卖
func (h *AuctionHandler) GetAuction(c *gin.Context) {
	idStr := c.Param("id")
	ctx := c.Request.Context()

	// 尝试按链上AuctionID查询
	if auctionID, err := strconv.ParseUint(idStr, 10, 64); err == nil {
		auction, err := h.service.GetAuctionByAuctionID(ctx, auctionID)
		if err != nil {
			// 尝试按数据库ID查询
			if dbID, err2 := strconv.ParseUint(idStr, 10, 32); err2 == nil {
				auction, err = h.service.GetAuctionByID(ctx, uint(dbID))
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{
						"success": false,
						"error":   "拍卖不存在",
					})
					return
				}
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"error":   "拍卖不存在",
				})
				return
			}
		}

		// 获取出价历史
		bids, _, err := h.service.GetAuctionBids(ctx, auction.AuctionID, 1, 20)
		if err != nil {
			// 即使获取出价历史失败，也返回拍卖信息
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"auction":     auction,
					"bid_history": []model.BidHistory{},
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"auction":     auction,
				"bid_history": bids,
			},
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "无效的拍卖ID",
	})
}

// CheckAuctionStatus 检查拍卖状态（只读查询）
func (h *AuctionHandler) CheckAuctionStatus(c *gin.Context) {
	identifier := c.Param("id")
	ctx := c.Request.Context()

	var auction *model.Auction
	var err error

	// 判断标识符类型
	if id, parseErr := strconv.ParseUint(identifier, 10, 64); parseErr == nil {
		// 链上AuctionID
		auction, err = h.service.GetAuctionByAuctionID(ctx, id)
		if err != nil {
			// 尝试数据库ID
			if dbID, parseErr2 := strconv.ParseUint(identifier, 10, 32); parseErr2 == nil {
				auction, err = h.service.GetAuctionByID(ctx, uint(dbID))
			}
		}
	}

	if err != nil || auction == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "未找到拍卖记录",
			"status":  "not_found",
		})
		return
	}

	// 计算剩余时间
	var timeRemaining uint64
	if auction.EndTime > uint64(time.Now().Unix()) {
		timeRemaining = auction.EndTime - uint64(time.Now().Unix())
	}

	// 构建响应
	response := gin.H{
		"success": true,
		"data": gin.H{
			"auction_id":     auction.AuctionID,
			"status":         auction.Status,
			"ended":          auction.Ended,
			"seller":         auction.Seller,
			"starting_price": auction.StartingPrice,
			"highest_bid":    auction.HighestBid,
			"highest_bidder": auction.HighestBidder,
			"start_time":     auction.StartTime,
			"end_time":       auction.EndTime,
			"time_remaining": timeRemaining,
			"nft_contract":   auction.NFTContract,
			"token_id":       auction.TokenID,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetAuctionBids 获取拍卖的出价历史
func (h *AuctionHandler) GetAuctionBids(c *gin.Context) {
	idStr := c.Param("id")
	ctx := c.Request.Context()

	auctionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的拍卖ID",
		})
		return
	}

	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	bids, total, err := h.service.GetAuctionBids(ctx, auctionID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取出价历史失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"bids": bids,
			"pagination": gin.H{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// ==================== 合约信息API ====================

// GetContractInfo 获取拍卖合约信息
func (h *AuctionHandler) GetContractInfo(c *gin.Context) {
	ctx := c.Request.Context()

	info, err := h.service.GetContractInfo(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取合约信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}

// GetAuctionCount 获取拍卖总数
func (h *AuctionHandler) GetAuctionCount(c *gin.Context) {
	ctx := c.Request.Context()

	count, err := h.service.GetAuctionCount(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取拍卖总数失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"auction_count": count,
		},
	})
}

// ==================== 验证API ====================

// ValidateAuction 验证拍卖是否存在
func (h *AuctionHandler) ValidateAuction(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	auctionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的拍卖ID格式",
		})
		return
	}

	exists, err := h.service.ValidateAuctionExists(ctx, auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "验证拍卖失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"auction_id": auctionID,
			"exists":     exists,
		},
	})
}

// ==================== 同步API（管理用）====================

// SyncAuctions 手动同步拍卖数据（管理用）
func (h *AuctionHandler) SyncAuctions(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.service.SyncAllAuctions(ctx); err != nil {
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
