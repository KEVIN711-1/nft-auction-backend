package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

// 用户浏览NFT市场
//
//	↓
//
// API网关 → 查询数据库（Redis缓存）← 返回数据（<100ms）
//
//	        ↑
//	监听服务（监听链上事件）
//	        ↑
//	   区块链节点
func main() {
	// ==================== 1. 配置加载阶段 ====================
	cfg := config.LoadConfig()
	log.SetPrefix("[NFT_BACK_END] ")

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
	// ┌─────────────┐    调用    ┌─────────────┐    调用    ┌─────────────┐
	// │   API层     │───────────▶│ Service层  │───────────▶│ Contract层  │
	// │  Handlers   │            │  Services  │            │    Client   │
	// └─────────────┘            └─────────────┘            └─────────────┘
	//        │                         │                            │
	//        │ 返回JSON                │ 业务逻辑                   │ 区块链交互
	//        ▼                         ▼                            ▼
	//    前端/客户端               数据库操作                    以太坊网络

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

	// 初始化拍卖客户端
	auctionClient, err := contract.NewAuctionClient(cfg.Blockchain.RPCURL, cfg.Blockchain.AuctionContractAddress)
	if err != nil {
		log.Fatalf("❌ 拍卖客户端初始化失败: %v", err)
	}

	// ==================== 4. 服务层初始化 ====================
	// user 服务
	userService := service.NewUserService(db)
	userHandler := api.NewUserHandler(userService)

	// NFT 服务
	nftService := service.NewNFTService(db, nftClient)
	nftHandler = api.NewNFTHandler(nftService)

	// NFT拍卖 服务（传入两个客户端）
	auctionService := service.NewAuctionService(db, auctionClient)
	auctionHandler := api.NewAuctionHandler(auctionService)

	// ==================== 5.区块链监听器初始化 ====================
	// 启动监听器（使用后台context）
	// context.WithCancel 是 Go 语言中用于创建 可取消的上下文（Context） 的函数
	//         case <-ctx.Done():  // 在监听器中 监听取消信号

	// 它在函数调用之间显式传递
	// 它携带本次调用的相关信息（取消信号、超时、请求ID等）
	// 每个请求/任务有自己独立的Context链条
	// 它让函数知道自己为什么运行、何时应该停止
	// 在你的区块链监听器中，ctx 让监听器知道："当主程序要退出时，请优雅地停止监听，清理资源"。
	log.SetPrefix("[NFT_LISTENER] ")

	ctx, cancel := context.WithCancel(context.Background())
	blockchainListener := service.NewBlockchainListener(
		nftService,            // NFT Service
		auctionService,        // Auction Service
		cfg.Blockchain.RPCURL, // RPC URL
		ctx,
		cancel,
	)
	defer cancel()

	blockchainListener.Start(ctx)
	// ==================== 6. Web服务器路由设置 ====================
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
	router := gin.Default()

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

	// ==================== 公开路由 ====================
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
		// listenerStatus := blockchainListener.GetStatus()
		// eventStats := blockchainListener.GetEventStats()

		c.JSON(200, gin.H{
			"service": "NFT Auction Marketplace",
			"version": "1.0.0",
			"features": gin.H{
				"blockchain_listener": true,
				"real_time_sync":      true,
				// "polling_interval":    blockchainListener.GetPollInterval().String(),
			},
			// "listener": listenerStatus,
			// "stats":    eventStats,
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

	auth := router.Group("/api")

	// Gin 参数获取方法的完整对比：
	// 方法	用途	示例	对应 AirPost 位置
	// c.Param("id")	路径参数（URL 路径中的变量）	/api/nfts/123 → "123"	URL 路径中
	// c.Query("id")	查询参数（URL ?后面的参数）	/api/nfts?id=123 → "123"	Params 标签页
	// c.PostForm("id")	表单参数（POST 表单数据）	id=123（表单提交）	Body (form-data)
	// c.GetHeader("X-ID")	请求头参数	X-ID: 123	Headers 标签页
	// c.ShouldBindJSON(&obj)	JSON 请求体	{"id": "123"}	Body (raw JSON)

	// ==================== 公开的拍卖查询API ====================
	// 根据交易哈希查询拍卖
	router.GET("/api/auctions/by-tx", auctionHandler.CheckAuctionStatus)

	// 检查拍卖状态（前端轮询）
	router.GET("/api/auctions/:id/status", auctionHandler.CheckAuctionStatus)

	// 拍卖列表和详情（公开）
	router.GET("/api/auctions", auctionHandler.GetAuctions)
	router.GET("/api/auctions/active", auctionHandler.GetActiveAuctions)
	router.GET("/api/auctions/count", auctionHandler.GetAuctionCount)
	router.GET("/api/auctions/:id", auctionHandler.GetAuction)
	router.GET("/api/auctions/:id/bids", auctionHandler.GetAuctionBids)
	router.GET("/api/auctions/:id/validate", auctionHandler.ValidateAuction)

	// NFT相关API（公开）
	router.GET("/api/nfts/:id", nftHandler.GetNFTInfo)
	router.GET("/api/nfts/:id/owner", nftHandler.GetNFTOwner)
	router.GET("/api/nfts/:id/validate/:address", nftHandler.ValidateOwnership)

	// ==================== 需要认证的API ====================
	auth.Use(authCheck) // 检查是否登录
	{
		// 用户相关API
		auth.GET("/user/profile", userHandler.GetProfile)

		// 管理API
		auth.POST("/auctions/sync", auctionHandler.SyncAuctions)
		auth.POST("/nft/sync", nftHandler.SyncNFTInfo)

		// 监听器控制API（需要认证）
		auth.POST("/listener/restart", func(c *gin.Context) {
			// 停止当前监听器
			blockchainListener.Stop()
			time.Sleep(1 * time.Second)

			// 重新启动
			blockchainListener.Start(ctx)

			c.JSON(200, gin.H{
				"success":   true,
				"message":   "区块链监听器已重启",
				"timestamp": time.Now().Unix(),
			})
		})

		auth.POST("/listener/force-sync", func(c *gin.Context) {
			// 强制全量同步
			go func() {
				if err := auctionService.SyncAllAuctions(ctx); err != nil {
					log.Printf("强制同步失败: %v", err)
				}
			}()

			c.JSON(200, gin.H{
				"success":   true,
				"message":   "已触发全量同步，请稍后查看结果",
				"timestamp": time.Now().Unix(),
			})
		})
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

	// 优雅关闭处理
	setupGracefulShutdown(cancel)

	// 启动HTTP服务器
	if err := router.Run(addr); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}

func setupGracefulShutdown(cancel context.CancelFunc) {
	// 监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		log.Printf("收到关闭信号: %v", sig)

		// 执行优雅关闭
		log.Println("正在停止区块链监听器...")
		cancel() // 这会触发监听器的停止

		// 等待一小段时间确保监听器完全停止
		time.Sleep(2 * time.Second)

		log.Println("系统已优雅关闭")
		os.Exit(0)
	}()
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
