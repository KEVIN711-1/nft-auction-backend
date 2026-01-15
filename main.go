package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"nft-auction-backend/api"               // API处理器层
	"nft-auction-backend/internal/config"   // 配置管理
	"nft-auction-backend/internal/contract" // 区块链交互层
	"nft-auction-backend/internal/service"  // 业务逻辑层
	"nft-auction-backend/pkg/database"      // 数据库层
)

// 全局token存储（添加互斥锁保证并发安全）
var (
	loginTokens = make(map[string]string) // token -> username
	tokenMutex  = &sync.RWMutex{}
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

	var nftHandler *api.NFTHandler

	// 检查必要的配置
	if cfg.Blockchain.RPCURL == "" {
		log.Fatal("❌ 请在 config.yaml 中配置 blockchain.rpc_url")
	}

	if cfg.Blockchain.NFTContractAddress == "" {
		log.Fatal("❌ 请在 config.yaml 中配置 blockchain.nft_contract_address")
	}

	// 初始化NFT客户端
	nftClient, err := contract.NewNFTClient(cfg.Blockchain.RPCURL, cfg.Blockchain.NFTContractAddress)
	if err != nil {
		log.Fatalf("❌ NFT客户端初始化失败: %v", err)
	}
	log.Println("✅ NFT客户端初始化成功")

	// 初始化拍卖客户端
	auctionClient, err := contract.NewAuctionClient(cfg.Blockchain.RPCURL, cfg.Blockchain.AuctionContractAddress)
	if err != nil {
		log.Fatalf("❌ 拍卖客户端初始化失败: %v", err)
	}
	log.Println("✅ 拍卖客户端初始化成功")

	// ==================== 4. 服务层初始化 ====================
	// 新增：用户服务初始化
	userService := service.NewUserService(db)
	// 新增：用户处理器
	userHandler := api.NewUserHandler(userService)

	// ┌─────────────┐    调用    ┌─────────────┐    调用    ┌─────────────┐
	// │   API层     │───────────▶│ Service层  │───────────▶│ Contract层  │
	// │  Handlers   │            │  Services  │            │    Client   │
	// └─────────────┘            └─────────────┘            └─────────────┘
	//        │                         │                            │
	//        │ 返回JSON                │ 业务逻辑                   │ 区块链交互
	//        ▼                         ▼                            ▼
	//    前端/客户端               数据库操作                    以太坊网络

	// ✅ NFT信息实时获取 → 因为所有权可能随时转移
	// ✅ 拍卖信息存数据库 → 因为需要历史记录和复杂查询
	// ✅ 混合架构 → 区块链应用的最佳实践
	// NFT服务
	nftService := service.NewNFTService(nftClient)
	nftHandler = api.NewNFTHandler(nftService)

	// 拍卖服务（传入两个客户端）
	auctionService := service.NewAuctionService(db, nftClient, auctionClient)
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
	// 	中间件（Middleware） = 在请求和响应之间的一系列处理函数
	// 特点：
	//     链式执行：一个接一个，像流水线
	//     可提前终止：任意环节可以"拦截"请求
	//     共享上下文：可以通过c.Set()/c.Get()传递数据
	//     顺序重要：先执行的中间件可能影响后续中间件

	// 你的CORS中间件在做什么？
	//     给每个响应"贴上标签"："允许跨域访问"
	//     专门处理浏览器"试探性"的OPTIONS请求
	//     让真正的业务逻辑（路由处理函数）不用关心跨域问题

	// 为什么叫"中间件"？
	// 因为它站在中间：
	//     不是客户端（浏览器）
	//     不是最终的业务逻辑
	//     是"中间的处理件"
	// 就相当于，提公因式，并且过滤一些不支持的请求或者放行一些特殊请求

	// 加上需要登录才能创建拍卖的中间件
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

	// ==================== 公开路由（新增：使用userHandler的方法） ====================
	router.POST("/register", userHandler.Register) // 注册 - 使用userHandler
	router.POST("/login", userHandler.Login)       // 登录 - 使用userHandler

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

	auth := router.Group("/api")

	// Gin 参数获取方法的完整对比：
	// 方法	用途	示例	对应 AirPost 位置
	// c.Param("id")	路径参数（URL 路径中的变量）	/api/nfts/123 → "123"	URL 路径中
	// c.Query("id")	查询参数（URL ?后面的参数）	/api/nfts?id=123 → "123"	Params 标签页
	// c.PostForm("id")	表单参数（POST 表单数据）	id=123（表单提交）	Body (form-data)
	// c.GetHeader("X-ID")	请求头参数	X-ID: 123	Headers 标签页
	// c.ShouldBindJSON(&obj)	JSON 请求体	{"id": "123"}	Body (raw JSON)

	// 拍卖相关API
	auth.Use(authCheck) // 检查是否登录
	{
		// 新增：用户相关API（需要认证）
		auth.GET("/user/profile", userHandler.GetProfile)

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
	}

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
	log.Println("  POST /api/auctions                  - 创建拍卖")   // ?
	log.Println("  POST /api/auctions/:id/bid          - 出价")     // ?
	log.Println("  POST /api/auctions/:id/end          - 结束拍卖")   // ?
	log.Println("  GET  /api/nfts/:id                  - NFT信息")  // ?
	log.Println("  GET  /api/nfts/:id/owner            - NFT所有者") // ?
	log.Println("  GET  /api/nfts/:id/validate/:addr   - 验证所有权")  // ?
	log.Println("  GET  /api/nfts/contract/info        - 获取合约信息") //?
	log.Println("========================================")

	// 启动HTTP服务器
	if err := router.Run(addr); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}

// 登录检查中间件（与你的博客系统一致）
func authCheck(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "请先登录"})
		c.Abort()
		return
	}

	// 检查token是否有效
	username, exists := loginTokens[token]
	if !exists {
		c.JSON(401, gin.H{"error": "登录已过期，请重新登录"})
		c.Abort()
		return
	}

	// 保存用户信息到上下文
	c.Set("username", username)
	c.Next()
}

func GenerateSimpleToken(username string) string {
	return fmt.Sprintf("%d-%s", time.Now().Unix(), username)
}
