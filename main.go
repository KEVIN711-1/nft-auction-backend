package main

import (
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
	// 对应文件: internal/config/config.go
	log.Println("🚀 启动NFT拍卖后端系统...")
	cfg := config.LoadConfig() // 从config.yaml加载所有配置 如合约地址、rpc_url 链接

	// ==================== 2. 数据库初始化阶段 ====================
	// 对应文件: pkg/database/gorm.go
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

	// ==================== 3. 区块链客户端初始化 ====================
	// 对应文件: internal/contract/contract.go (基础客户端)
	log.Println("正在初始化区块链客户端...")

	// 检查RPC URL配置
	if cfg.Blockchain.RPCURL == "" {
		log.Println("⚠️  配置文件中rpc_url为空，请检查config.yaml")
		log.Println("📡 使用模拟模式运行")
	}

	// 创建基础区块链客户端（用于拍卖合约）
	contractClient, err := contract.NewContractClient(cfg.Blockchain.RPCURL)
	if err != nil {
		log.Printf("⚠️  区块链客户端初始化失败: %v", err)
		log.Println("📡 使用模拟模式运行")
		contractClient, _ = contract.NewContractClient("")
	}

	if contractClient.IsActive() {
		log.Println("✅ 区块链客户端连接成功")
	} else {
		log.Println("📡 运行在演示模式（使用模拟数据）")
	}

	// ==================== 4. NFT客户端和服务初始化 ====================
	// 对应文件:
	//   - internal/contract/nft_client.go (NFT专用客户端)
	//   - internal/service/nft_service.go (NFT业务逻辑)
	//   - api/nft.go (NFT API处理器)
	var nftHandler *api.NFTHandler

	// 创建NFT专用客户端（连接到你的NFT合约）
	nftClient, err := contract.NewNFTClient(cfg.Blockchain.RPCURL, cfg.Blockchain.NFTContractAddress)
	if err != nil {
		log.Printf("⚠️  NFT客户端初始化失败: %v", err)
		log.Println("📡 将继续运行，但无法获取NFT信息")
	} else {
		log.Println("✅ NFT客户端初始化成功")

		// 初始化NFT业务服务
		nftService := service.NewNFTService(db, nftClient)

		// 初始化NFT API处理器
		nftHandler = api.NewNFTHandler(nftService)

		// 首次同步NFT信息（异步）
		// 为什么异步？避免阻塞主线程，让服务器快速启动
		go func() {
			time.Sleep(3 * time.Second) // 等待其他服务初始化完成
			log.Println("🔄 开始同步NFT信息...")
			if err := nftService.SyncNFTInfo(); err != nil {
				log.Printf("首次NFT信息同步失败: %v", err)
			} else {
				log.Println("✅ 首次NFT信息同步完成")
			}
		}()
	}

	// ==================== 5. 拍卖服务初始化 ====================
	// 对应文件:
	//   - internal/service/auction_service.go (拍卖业务逻辑)
	//   - api/auction.go (拍卖API处理器)
	auctionService := service.NewAuctionService(db, contractClient)

	// ==================== 6. 异步任务启动 ====================
	// 首次同步拍卖数据（异步）
	go func() {
		time.Sleep(2 * time.Second) // 等待服务完全启动
		log.Println("🔄 开始首次数据同步...")
		if err := auctionService.SyncAuctions(); err != nil {
			log.Printf("首次同步失败: %v", err)
		} else {
			log.Println("✅ 首次数据同步完成")
		}
	}()

	// ==================== 7. API处理器初始化 ====================
	auctionHandler := api.NewAuctionHandler(auctionService)

	// ==================== 8. Web服务器路由设置 ====================
	// 使用Gin框架创建HTTP服务器
	router := gin.Default()

	// CORS中间件 - 允许前端跨域访问
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
	// 格式: router.HTTP方法("路径", 处理函数)

	// 9.1 健康检查API - 简单状态检查
	// 调用链: 前端 → Gin → 匿名函数
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"service":   "NFT Auction Backend",
		})
	})

	// 9.2 系统信息API - 返回当前配置信息
	// 调用链: 前端 → Gin → 匿名函数
	router.GET("/api/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "NFT Auction Marketplace",
			"version": "1.0.0",
			"config": gin.H{
				"port":             cfg.Server.Port,
				"database":         cfg.Database.Path,
				"rpc_url":          cfg.Blockchain.RPCURL,
				"nft_contract":     cfg.Blockchain.NFTContractAddress,
				"auction_contract": cfg.Blockchain.AuctionContractAddress,
			},
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// 9.3 拍卖相关API
	// 调用链: 前端 → Gin → auction_handler → auction_service → database/contract
	router.GET("/api/auctions", auctionHandler.GetAuctions)              // 获取所有拍卖
	router.GET("/api/auctions/active", auctionHandler.GetActiveAuctions) // 获取进行中拍卖
	router.GET("/api/auctions/:id", auctionHandler.GetAuction)           // 获取单个拍卖详情
	router.POST("/api/auctions/sync", auctionHandler.SyncAuctions)       // 手动同步拍卖数据

	// 9.4 NFT相关API
	// 调用链: 前端 → Gin → nft_handler → nft_service → nft_client → 区块链
	if nftHandler != nil {
		router.GET("/api/nft/info", nftHandler.GetNFTInfo)   // 获取NFT信息
		router.POST("/api/nft/sync", nftHandler.SyncNFTInfo) // 手动同步NFT信息
	}

	// ==================== 10. 服务器启动 ====================
	port := cfg.Server.Port
	addr := fmt.Sprintf(":%d", port)

	// 打印启动信息
	log.Println("========================================")
	log.Println("🎉 NFT拍卖后端系统启动成功!")
	log.Printf("📡 服务地址: http://localhost:%d", port)
	log.Printf("💾 数据库文件: %s", cfg.Database.Path)
	if cfg.Blockchain.RPCURL != "" {
		log.Printf("🔗 区块链节点: %s", cfg.Blockchain.RPCURL)
	}
	log.Println("========================================")
	log.Println("🌐 可用API接口:")
	log.Println("  GET  /api/health          - 健康检查")
	log.Println("  GET  /api/info            - 系统信息")
	log.Println("  GET  /api/auctions        - 所有拍卖")
	log.Println("  GET  /api/auctions/active - 进行中拍卖")
	log.Println("  GET  /api/auctions/:id    - 单个拍卖详情")
	log.Println("  POST /api/auctions/sync   - 手动同步数据")
	if nftHandler != nil {
		log.Println("  GET  /api/nft/info        - NFT信息")
		log.Println("  POST /api/nft/sync        - 同步NFT信息")
	}
	log.Println("========================================")

	// 启动HTTP服务器
	if err := router.Run(addr); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}
