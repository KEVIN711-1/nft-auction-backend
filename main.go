package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"nft-auction-backend/api"               // API处理器层
	"nft-auction-backend/internal/config"   // 配置管理
	"nft-auction-backend/internal/contract" // 区块链交互层
	"nft-auction-backend/internal/service"  // 业务逻辑层
	"nft-auction-backend/pkg/database"      // 数据库层
)

func main() {
	// ==================== 1. 配置加载阶段 ====================
	log.Println("🚀 启动NFT拍卖后端系统...")
	cfg := config.LoadConfig()

	// ==================== 2. 数据库初始化阶段 ====================
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// ==================== 3. NFT客户端初始化 ====================
	log.Println("正在初始化NFT客户端...")

	var nftClient contract.NFTContract
	var nftHandler *api.NFTHandler

	// 检查必要的配置
	if cfg.Blockchain.RPCURL == "" {
		log.Fatal("❌ 请在 config.yaml 中配置 blockchain.rpc_url")
	}

	if cfg.Blockchain.NFTContractAddress == "" {
		log.Fatal("❌ 请在 config.yaml 中配置 blockchain.nft_contract_address")
	}

	// 创建NFT客户端
	nftClient, err = contract.NewNFTClient(cfg.Blockchain.RPCURL, cfg.Blockchain.NFTContractAddress)
	if err != nil {
		log.Fatalf("❌ NFT客户端初始化失败: %v", err)
	}

	log.Println("✅ NFT客户端初始化成功")

	// ==================== 4. 服务层初始化 ====================
	// NFT服务
	nftService := service.NewNFTService(nftClient)
	nftHandler = api.NewNFTHandler(nftService)

	// 拍卖服务
	auctionService := service.NewAuctionService(db, nftClient)
	auctionHandler := api.NewAuctionHandler(auctionService)

	// ==================== 5. 测试连接 ====================
	go func() {
		time.Sleep(2 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 测试获取合约信息
		name, err := nftClient.GetName(ctx)
		if err != nil {
			log.Printf("⚠️  测试连接失败 - 无法获取合约名称: %v", err)
		} else {
			log.Printf("✅ 合约连接正常 - 名称: %s", name)

			// 测试获取symbol
			symbol, err := nftClient.GetSymbol(ctx)
			if err != nil {
				log.Printf("⚠️  无法获取合约符号: %v", err)
			} else {
				log.Printf("✅ 合约符号: %s", symbol)
			}
		}
	}()

	// ==================== 6. Web服务器路由设置 ====================
	router := gin.Default()

	// CORS中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	// ==================== API路由注册 ====================
	// 健康检查
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"service":   "NFT Auction Backend",
		})
	})

	// 系统信息
	router.GET("/api/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "NFT Auction Marketplace",
			"version": "1.0.0",
			"config": gin.H{
				"port":         cfg.Server.Port,
				"database":     cfg.Database.Path,
				"rpc_url":      cfg.Blockchain.RPCURL,
				"nft_contract": cfg.Blockchain.NFTContractAddress,
			},
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// 拍卖相关API
	router.GET("/api/auctions", auctionHandler.GetAuctions)
	router.GET("/api/auctions/active", auctionHandler.GetActiveAuctions)
	router.GET("/api/auctions/:id", auctionHandler.GetAuction)
	router.POST("/api/auctions", auctionHandler.CreateAuction)
	router.POST("/api/auctions/:id/bid", auctionHandler.PlaceBid)
	router.POST("/api/auctions/:id/end", auctionHandler.EndAuction)
	router.POST("/api/auctions/sync", auctionHandler.SyncAuctions)

	// NFT相关API
	router.GET("/api/nfts/:id", nftHandler.GetNFTInfo)
	router.GET("/api/nfts/:id/owner", nftHandler.GetNFTOwner)
	router.GET("/api/nfts/:id/validate/:address", nftHandler.ValidateOwnership)
	router.GET("/api/nfts/contract/info", nftHandler.GetContractInfo)
	router.POST("/api/nft/sync", nftHandler.SyncNFTInfo)

	// ==================== 7. 服务器启动 ====================
	port := cfg.Server.Port
	addr := fmt.Sprintf(":%d", port)

	// 打印启动信息
	log.Println("========================================")
	log.Println("🎉 NFT拍卖后端系统启动成功!")
	log.Printf("📡 服务地址: http://localhost:%d", port)
	log.Printf("💾 数据库文件: %s", cfg.Database.Path)
	log.Printf("🔗 区块链节点: %s", cfg.Blockchain.RPCURL)
	log.Printf("📄 NFT合约地址: %s", cfg.Blockchain.NFTContractAddress)
	log.Println("========================================")
	log.Println("🌐 可用API接口:")
	log.Println("  GET  /api/health                    - 健康检查")
	log.Println("  GET  /api/info                      - 系统信息")
	log.Println("  GET  /api/auctions                  - 所有拍卖")
	log.Println("  GET  /api/auctions/active           - 进行中拍卖")
	log.Println("  GET  /api/auctions/:id              - 单个拍卖详情")
	log.Println("  POST /api/auctions                  - 创建拍卖")
	log.Println("  POST /api/auctions/:id/bid          - 出价")
	log.Println("  POST /api/auctions/:id/end          - 结束拍卖")
	log.Println("  GET  /api/nfts/:id                  - NFT信息")
	log.Println("  GET  /api/nfts/:id/owner            - NFT所有者")
	log.Println("  GET  /api/nfts/:id/validate/:addr   - 验证所有权")
	log.Println("  GET  /api/nfts/contract/info        - 获取合约信息")
	log.Println("========================================")

	// 启动HTTP服务器
	if err := router.Run(addr); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}
