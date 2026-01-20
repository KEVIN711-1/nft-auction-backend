package service

import (
	"context"
	"log"
	"math/big"
	"sync"
	"time"

	"nft-auction-backend/internal/contract"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ────────────────────┐
// │    区块链网络       │
// │  (NFT & Auction)   │
// └─────────┬──────────┘
//           │ 链上事件（Transfer, AuctionCreated, BidPlaced, AuctionEnded）
//           ▼
// ┌────────────────────┐
// │  区块链监听模块     │  <- BlockchainListener
// │  - WebSocket/RPC   │
// │  - 解析事件         │
// │  - 去重/校验       │
// └─────────┬──────────┘
//           │ 解析后的结构化数据
//           ▼
// ┌────────────────────┐
// │     后端数据库      │  <- NFTInfo, Auction, Bid 表
// │  - 更新 NFT 所有权  │
// │  - 保存出价记录     │
// │  - 更新拍卖状态     │
// └─────────┬──────────┘
//           │ 提供接口/触发通知
//           ▼
// ┌────────────────────┐
// │ 后端 API / WebSocket│
// │  - REST API         │
// │    GET /nfts/:id    │
// │    GET /auctions    │
// │  - WebSocket/SSE    │
// │    实时推送事件     │
// └─────────┬──────────┘
//           │ JSON 数据 / 实时事件
//           ▼
// ┌────────────────────┐
// │      前端页面       │
// │  - NFT 拥有者显示   │
// │  - 最新出价显示     │
// │  - 拍卖状态更新     │
// └────────────────────┘

// BlockchainListener 监听区块链事件
type BlockchainListener struct {
	rpcURL         string
	ethClient      *ethclient.Client
	nftService     *NFTService
	auctionService *AuctionService

	ctx       context.Context
	cancel    context.CancelFunc
	running   bool
	stats     map[string]int
	statsLock sync.RWMutex
}

// NewBlockchainListener 创建监听器
func NewBlockchainListener(
	nft contract.NFTContract,
	auction contract.AuctionContract,
	nftSvc *NFTService,
	auctionSvc *AuctionService,
	rpcURL string,
) *BlockchainListener {
	ctx, cancel := context.WithCancel(context.Background())
	return &BlockchainListener{
		rpcURL:         rpcURL,
		nftService:     nftSvc,
		auctionService: auctionSvc,
		ctx:            ctx,
		cancel:         cancel,
		stats:          map[string]int{"nft_transfers": 0, "auctions": 0, "bids": 0},
	}
}

// Start 启动监听器
func (l *BlockchainListener) Start(ctx context.Context) {
	if l.running {
		return
	}
	l.running = true
	log.Println("🔍 区块链事件监听器启动中...")

	go func() {
		// 无限循环，持续监听区块链事件
		// 除非收到停止信号（ctx.Done()），否则会一直运行
		for {
			select {
			case <-l.ctx.Done():
				// 如果收到停止信号，输出日志并退出函数
				log.Println("❌ 区块链监听器已停止")
				return
			default:
				// 连接 WebSocket RPC
				log.Println("----1----🔄 区块链监听器开始同步...")

				client, err := ethclient.Dial(l.rpcURL)
				if err != nil {
					log.Printf("❌ 连接 RPC 失败: %v, 3s后重试...", err)
					time.Sleep(3 * time.Second)
					continue
				}
				l.ethClient = client

				// 1️⃣ 启动链上数据同步
				l.syncAllNFTs()
				l.syncAllAuctions()

				// 创建 WaitGroup，用于等待两个监听goroutine完成
				// 2️⃣ 启动 NFT 和拍卖实时监听
				var wg sync.WaitGroup

				// 设置需要等待的 goroutine 数量为 2
				wg.Add(2)
				go func() {
					defer wg.Done()
					l.listenNFTTransfer()
				}()
				go func() {
					defer wg.Done()
					l.listenAuctionEvents()
				}()

				// 等待两个监听任务完成
				// wg.Wait() 会阻塞，直到两个任务都调用了 wg.Done()
				// 这意味着只有当两个监听函数都退出时，才会继续执行后面的代码
				wg.Wait()

				// 如果监听退出，关闭客户端重连
				l.ethClient.Close()
				log.Println("----2----🔄 区块链监听器重连中...")
				time.Sleep(3 * time.Second)
			}
		}
	}()
}

// Stop 停止监听器
func (l *BlockchainListener) Stop() {
	if !l.running {
		return
	}
	log.Println("🛑 停止区块链监听器...")
	l.cancel()
	if l.ethClient != nil {
		l.ethClient.Close()
	}
	l.running = false
}

// ---------------- 拍卖同步 ----------------
func (l *BlockchainListener) syncAllAuctions() {
	log.Println("⏳ 同步链上所有拍卖数据中...")
	if err := l.auctionService.SyncAllAuctions(l.ctx); err != nil {
		log.Printf("❌ 同步拍卖失败: %v", err)
		return
	}
	log.Println("✅ 拍卖同步完成")
}

func (l *BlockchainListener) syncAllNFTs() {
	log.Println("⏳ 同步链上所有NFT数据中...")
	if err := l.nftService.SyncAllNFTs(l.ctx); err != nil {
		log.Printf("❌ 同步拍卖失败: %v", err)
		return
	}
	log.Println("✅ 拍卖同步完成")
}

// ---------------- NFT Transfer 监听 ----------------
func (l *BlockchainListener) listenNFTTransfer() {
	nftAddr := l.nftService.GetContractAddress()
	log.Printf("🎯 监听合约: %s", nftAddr.Hex())
	query := ethereum.FilterQuery{Addresses: []common.Address{nftAddr}}
	logsChan := make(chan types.Log)

	// 计算预期的签名
	transferSig := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Hex()
	mintSig := crypto.Keccak256Hash([]byte("NFTMinted(address,uint256,string)")).Hex()

	log.Printf("  Transfer签名: %s", transferSig)
	log.Printf("  Minted签名: %s", mintSig)
	// 从最新区块开始监听
	latestBlock, err := l.ethClient.BlockNumber(l.ctx)
	if err == nil {
		log.Printf("📦 从区块 #%d 开始监听", latestBlock)
	}

	// SubscribeFilterLogs默认从最新区块开始监听
	sub, err := l.ethClient.SubscribeFilterLogs(l.ctx, query, logsChan)
	if err != nil {
		log.Fatalf("❌ 订阅失败: %v", err)
	}
	log.Println("✅1 NFT 事件监听器订阅成功，等待事件...")

	for {
		select {
		case err := <-sub.Err():
			log.Printf("❌ 订阅错误: %v", err)
			return

		case vLog := <-logsChan:
			if len(vLog.Topics) == 0 {
				continue
			}
			// 打印事件基本信息
			log.Printf("📥 NFT  监听器收到事件日志:")
			log.Printf("  区块: %d", vLog.BlockNumber)
			log.Printf("  交易: %s", vLog.TxHash.Hex())
			log.Printf("  主题数: %d", len(vLog.Topics))
			// 监听到的事件签名
			eventSig := vLog.Topics[0].Hex()
			log.Printf("  事件签名: %s", eventSig)

			if eventSig == mintSig {
				mintEvent, err := l.nftService.client.ParseNFTMinted(vLog)
				if err == nil {
					log.Printf("✅ 解析到Mint事件: TokenID=%s", mintEvent.TokenId)
					// 理想状态下为获取事件传递的参数后，只更新参数，不用大费周章再根据id 去拉取一边区块链的信息了
					err := l.nftService.UpdateNFTFromChain(mintEvent.TokenId.String())
					if err != nil {
						log.Printf("❌ NFT同步失败: %v", err)
						continue
					}
					log.Printf("✅ NFT已同步: TokenID=%s", mintEvent.TokenId.String())
					continue
				} else {
					log.Printf("❌ 解析Mint事件失败: %v", err)
				}
			} else if eventSig == transferSig {
				// 尝试解析Transfer事件
				transferEvent, err := l.nftService.client.ParseTransfer(vLog)
				if err == nil {
					log.Printf("✅ 解析到Transfer事件: TokenID=%s", transferEvent.TokenId)
					// 理想状态下为获取事件传递的参数后，只更新参数，不用大费周章再根据id 去拉取一边区块链的信息了
					err := l.nftService.UpdateNFTFromChain(transferEvent.TokenId.String())
					if err != nil {
						log.Printf("❌ NFT同步失败: %v", err)
						continue
					}
					log.Printf("✅ NFT已同步: TokenID=%s", transferEvent.TokenId.String())
					continue
				} else {
					log.Printf("❌ 解析Transfer事件失败: %v", err)
				}
			} else {
				// 加一些approve 的监听
				log.Printf("⚠️ 无法解析的事件，跳过")
			}
		case <-l.ctx.Done():
			log.Println("🛑 监听器停止")
			return
		}
	}
}

// ---------------- 拍卖事件监听 ----------------
func (l *BlockchainListener) listenAuctionEvents() {
	auctionAddr := l.auctionService.GetContractAddress()
	query := ethereum.FilterQuery{Addresses: []common.Address{auctionAddr}}

	// 根据你的合约声明，正确的签名计算：
	// 注意：参数顺序和类型必须完全匹配
	auctionCreatedID := crypto.Keccak256Hash([]byte("AuctionCreated(uint256,address,uint256,uint256)"))
	bidPlacedID := crypto.Keccak256Hash([]byte("NewBid(uint256,address,uint256)"))
	auctionEndedID := crypto.Keccak256Hash([]byte("AuctionEnded(uint256,address,uint256)"))

	// 调试输出
	log.Printf("📊 计算的事件签名:")
	log.Printf("  AuctionCreated: %s", auctionCreatedID.Hex())
	log.Printf("  NewBid: %s", bidPlacedID.Hex())
	log.Printf("  AuctionEnded: %s", auctionEndedID.Hex())

	log.Println("✅2 NFT 拍卖事件监听器订阅成功，等待事件...")
	logsChan := make(chan types.Log)
	sub, err := l.ethClient.SubscribeFilterLogs(l.ctx, query, logsChan)
	if err != nil {
		log.Fatalf("❌ 拍卖 SubscribeFilterLogs 失败: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Printf("❌ 拍卖监听错误: %v, 重连中...", err)
			return
		case vLog := <-logsChan:
			if len(vLog.Topics) == 0 {
				continue
			}

			eventID := vLog.Topics[0]
			log.Printf("📥 NFT 拍卖事件收到事件，签名: %s", eventID.Hex())

			// 重要：你的事件参数都不是indexed，所以auctionId在Data字段，不在Topics中
			var auctionID *big.Int

			switch eventID {
			case auctionCreatedID, bidPlacedID, auctionEndedID:
				// 因为参数没有indexed，auctionId在Data字段的前32字节
				if len(vLog.Data) >= 32 {
					auctionID = new(big.Int).SetBytes(vLog.Data[:32])
				}

				var name string
				switch eventID {
				case auctionCreatedID:
					name = "AuctionCreated"
				case bidPlacedID:
					name = "NewBid"
				case auctionEndedID:
					name = "AuctionEnded"
				}

				if auctionID != nil {
					log.Printf("🏷️ 拍卖事件: %s, AuctionID=%s", name, auctionID.String())

					l.statsLock.Lock()
					l.stats["auctions"]++
					if eventID == bidPlacedID {
						l.stats["bids"]++
					}
					l.statsLock.Unlock()
					// 理想状态下为获取事件传递的参数后，只更新参数，不用大费周章再根据id 去拉取一边区块链的信息了
					if err := l.auctionService.UpdateAuctionFromChain(auctionID.Uint64()); err != nil {
						log.Printf("❌ 更新拍卖失败: %v", err)
					}
				}
			default:
				log.Printf("⚠️ 未知拍卖事件: %s", eventID.Hex())
			}
		case <-l.ctx.Done():
			log.Println("❌ 拍卖监听器已停止")
			return
		}
	}
}
