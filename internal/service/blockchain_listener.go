package service

import (
	"bytes"
	"context"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"nft-auction-backend/internal/contract"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	rpcURL          string
	ethClient       *ethclient.Client
	nftContract     contract.NFTContract
	auctionContract contract.AuctionContract

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
		rpcURL:          rpcURL,
		nftContract:     nft,
		auctionContract: auction,
		nftService:      nftSvc,
		auctionService:  auctionSvc,
		ctx:             ctx,
		cancel:          cancel,
		stats:           map[string]int{"nft_transfers": 0, "auctions": 0, "bids": 0},
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
	log.Println("⏳ 同步链上所有拍卖数据中...")
	if err := l.nftService.SyncAllNFTs(l.ctx); err != nil {
		log.Printf("❌ 同步拍卖失败: %v", err)
		return
	}
	log.Println("✅ 拍卖同步完成")
}

// ---------------- NFT Transfer 监听 ----------------
func (l *BlockchainListener) listenNFTTransfer() {
	nftAddr := l.nftContract.GetContractAddress()
	query := ethereum.FilterQuery{Addresses: []common.Address{nftAddr}}

	data, err := os.ReadFile("./internal/contract/abi/abi/KevinNFT.abi")
	if err != nil {
		log.Fatalf("❌ 读取 NFT ABI 文件失败: %v", err)
	}
	parsedABI, err := abi.JSON(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("❌ 解析 NFT ABI 失败: %v", err)
	}

	transferEvent := parsedABI.Events["Transfer"]

	log.Println("🔔 NFT Transfer 监听器已启动")
	logsChan := make(chan types.Log)
	sub, err := l.ethClient.SubscribeFilterLogs(l.ctx, query, logsChan)
	if err != nil {
		log.Fatalf("❌ NFT SubscribeFilterLogs 失败: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Printf("❌ NFT监听错误: %v, 重连中...", err)
			return
		case vLog := <-logsChan:
			if len(vLog.Topics) == 0 {
				continue
			}
			if vLog.Topics[0] == transferEvent.ID {
				tokenID := new(big.Int).SetBytes(vLog.Data)
				log.Printf("🔄 NFT Transfer 事件: TokenID=%s", tokenID.String())

				l.statsLock.Lock()
				l.stats["nft_transfers"]++
				l.statsLock.Unlock()

				err := l.nftService.UpdateNFTFromChain(tokenID.String())
				if err != nil {
					log.Printf("❌ NFT同步失败: %v", err)
					continue
				}
				log.Printf("✅ NFT已同步: TokenID=%s", tokenID.String())
			}
		case <-l.ctx.Done():
			log.Println("❌ NFT监听器已停止")
			return
		}
	}
}

// ---------------- 拍卖事件监听 ----------------
func (l *BlockchainListener) listenAuctionEvents() {
	auctionAddr := l.auctionContract.GetContractAddress()
	query := ethereum.FilterQuery{Addresses: []common.Address{auctionAddr}}

	data, err := os.ReadFile("./internal/contract/abi/abi/NftAuction.abi")
	if err != nil {
		log.Fatalf("❌ 读取拍卖 ABI 文件失败: %v", err)
	}
	parsedABI, err := abi.JSON(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("❌ 解析拍卖 ABI 失败: %v", err)
	}

	auctionCreatedID := parsedABI.Events["AuctionCreated"].ID
	bidPlacedID := parsedABI.Events["NewBid"].ID
	auctionEndedID := parsedABI.Events["AuctionEnded"].ID

	log.Println("🔔 拍卖事件监听器已启动")
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
			var auctionID *big.Int

			switch eventID {
			case auctionCreatedID, bidPlacedID, auctionEndedID:
				auctionID = new(big.Int).SetBytes(vLog.Data)
				var name string
				switch eventID {
				case auctionCreatedID:
					name = "AuctionCreated"
				case bidPlacedID:
					name = "BidPlaced"
				case auctionEndedID:
					name = "AuctionEnded"
				}
				log.Printf("🏷️ 拍卖事件: %s, AuctionID=%s", name, auctionID.String())

				l.statsLock.Lock()
				l.stats["auctions"]++
				if eventID == bidPlacedID {
					l.stats["bids"]++
				}
				l.statsLock.Unlock()

				if err := l.auctionService.UpdateAuctionFromChain(auctionID.Uint64()); err != nil {
					log.Printf("❌ 更新拍卖失败: %v", err)
				}
			}
		case <-l.ctx.Done():
			log.Println("❌ 拍卖监听器已停止")
			return
		}
	}
}
