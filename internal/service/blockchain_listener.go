package service

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"nft-auction-backend/internal/contract"
	"nft-auction-backend/internal/model"
	"sync"
	"time"

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
	nftSvc *NFTService,
	auctionSvc *AuctionService,
	rpcURL string,
	ctx context.Context,
	cancel context.CancelFunc) *BlockchainListener {

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
				log.Println(" 区块链监听器开始同步...")

				client, err := ethclient.Dial(l.rpcURL)
				if err != nil {
					log.Printf("❌ 连接 RPC 失败: %v, 3s后重试...", err)
					time.Sleep(3 * time.Second)
					continue
				}
				defer client.Close()
				l.ethClient = client

				// 先同步一遍链上的数据
				l.syncAllNFTs()
				l.syncAllAuctions()

				// 创建 WaitGroup，用于等待两个监听goroutine完成
				// 启动 NFT 和拍卖实时监听
				var wg sync.WaitGroup

				// 设置需要等待的 goroutine 数量为 2
				wg.Add(2)
				go func() {
					// 监听NFT 的监听器
					defer wg.Done() // 无论函数如何结束，defer都会执行
					l.listenNFTTransfer()
				}()
				go func() {
					// 监听NFT拍卖 的监听器
					defer wg.Done()
					l.listenAuctionEvents()
				}()

				// 等待两个监听任务完成
				// wg.Wait() 会阻塞，直到两个任务都调用了 wg.Done()
				// 这意味着只有当两个监听函数都退出时，才会继续执行后面的代码
				wg.Wait()
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
	log.Println("====1====⏳ 同步链上所有拍卖数据中...")
	if err := l.auctionService.SyncAllAuctions(l.ctx); err != nil {
		log.Printf("❌ 同步拍卖失败: %v", err)
		return
	}
}

func (l *BlockchainListener) syncAllNFTs() {
	log.Println("====2====⏳ 同步链上所有NFT数据中...")
	if err := l.nftService.SyncAllNFTs(l.ctx); err != nil {
		log.Printf("❌ 同步拍卖失败: %v", err)
		return
	}
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
	approvalSig := crypto.Keccak256Hash([]byte("Approval(address,address,uint256)")).Hex()

	log.Printf("  Transfer签名: %s", transferSig)
	log.Printf("  Minted签名: %s", mintSig)
	filterer, err := contract.NewKevinNFTFilterer(nftAddr, nil)
	if err != nil {
		log.Fatalf("❌ 创建Filterer失败: %v", err)
	}

	// SubscribeFilterLogs默认从最新区块开始监听
	sub, err := l.ethClient.SubscribeFilterLogs(l.ctx, query, logsChan)
	if err != nil {
		log.Fatalf("❌ 订阅失败: %v", err)
	}
	log.Println("✅ 1 NFT 事件监听器订阅成功，等待事件...")

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

			// 根据事件签名分流处理
			switch eventSig {
			case mintSig:
				l.handleNFTMinted(vLog, filterer)
			case transferSig:
				l.handleTransfer(vLog, filterer)
			case approvalSig:
				l.handleApproval(vLog, filterer)
			default:
				log.Printf("⚠️ 未知NFT事件签名: %s", eventSig)
			}
		case <-l.ctx.Done():
			log.Println("🛑 监听器停止")
			return
		}
	}
}

// ==================== 事件处理函数 ====================
// handleNFTMinted 处理NFT铸造事件
func (l *BlockchainListener) handleNFTMinted(vLog types.Log, filterer *contract.KevinNFTFilterer) {
	event, err := filterer.ParseNFTMinted(vLog)
	if err != nil {
		log.Printf("❌ 解析Mint事件失败: %v", err)
		return
	}

	log.Printf("✅ Mint事件: TokenID=%s, Owner=%s, URI=%s",
		event.TokenId.String(), event.Owner.Hex(), event.Uri)
	contractName, _ := l.nftService.client.GetName(l.ctx)
	contractSymbol, _ := l.nftService.client.GetSymbol(l.ctx)

	// 获取总供应量
	var totalSupply string
	if total, err := l.nftService.client.GetTotalSupply(l.ctx); err == nil {
		totalSupply = total.String()
	}
	// 直接从事件数据创建NFT记录，不需要再查询区块链
	nft := &model.NFTInfo{
		ContractAddress: l.nftService.GetContractAddress().Hex(),
		TokenID:         event.TokenId.String(),
		Owner:           event.Owner.Hex(),
		Name:            fmt.Sprintf("NFT #%s", event.TokenId.String()),
		Uri:             event.Uri,
		TotalSupply:     totalSupply,
		Blockchain:      "sepolia",
		ContractName:    contractName,
		ContractSymbol:  contractSymbol,
		IsMinted:        true,
		LastSyncTime:    time.Now(),
	}

	if err := l.nftService.SaveNFT(l.ctx, nft); err != nil {
		log.Printf("❌ 保存NFT失败: %v", err)
	} else {
		log.Printf("✅ NFT已保存: TokenID=%s", event.TokenId.String())
	}
}

// handleTransfer 处理NFT转移事件
func (l *BlockchainListener) handleTransfer(vLog types.Log, filterer *contract.KevinNFTFilterer) {
	event, err := filterer.ParseTransfer(vLog)
	if err != nil {
		log.Printf("❌ 解析Transfer事件失败: %v", err)
		return
	}

	log.Printf("✅ Transfer事件: TokenID=%s, From=%s, To=%s",
		event.TokenId.String(), event.From.Hex(), event.To.Hex())

	// 直接更新NFT所有者，不需要查询区块链
	contractAddr := l.nftService.GetContractAddress().Hex()
	tokenID := event.TokenId.String()
	newOwner := event.To.Hex()

	var existing model.NFTInfo
	result := l.nftService.DB.WithContext(l.ctx).
		Model(&model.NFTInfo{}).
		Where("contract_address = ? AND token_id = ?", contractAddr, tokenID).First(&existing)
	if result.Error != nil {
		log.Printf("❌ 数据库更新失败: %v", err)
	}
	existing.Owner = newOwner

	// 更新数据库中的NFT所有者
	if err := l.nftService.SaveNFT(l.ctx, &existing); err != nil {
		log.Printf("❌ 保存NFT失败: %v", err)
	} else {
		log.Printf("✅ NFT已保存: TokenID=%s", event.TokenId.String())
	}
}

// handleApproval 处理单NFT授权事件
func (l *BlockchainListener) handleApproval(vLog types.Log, filterer *contract.KevinNFTFilterer) {
	event, err := filterer.ParseApproval(vLog)
	if err != nil {
		log.Printf("❌ 解析Approval事件失败: %v", err)
		return
	}

	log.Printf("✅ Approval事件: TokenID=%s, Owner=%s, Approved=%s",
		event.TokenId.String(), event.Owner.Hex(), event.Approved.Hex())

	// 保存授权记录到数据库
	approval := &model.NFTInfo{
		TokenID:         event.TokenId.String(),
		Owner:           event.Owner.Hex(),
		ApprovedAddress: event.Approved.Hex(),
		ApprovedAt:      time.Now(),
		ApprovalTxHash:  vLog.TxHash.Hex(),
		LastSyncTime:    time.Now(),
	}

	if err := l.nftService.SaveNFT(l.ctx, approval); err != nil {
		log.Printf("❌ 保存授权记录失败: %v", err)
	}
}

// ==================== 辅助函数 ====================

// createNFTFromTransfer 从转移事件创建NFT记录
func (l *BlockchainListener) createNFTFromTransfer(contractAddr, tokenID, owner string) error {
	// 这里可以添加一些默认值或从区块链获取基本信息
	nft := &model.NFTInfo{
		ContractAddress: contractAddr,
		TokenID:         tokenID,
		Owner:           owner,
		Name:            fmt.Sprintf("NFT #%s", tokenID),
		Uri:             "", // 可能需要查询
		Blockchain:      "sepolia",
		IsMinted:        true,
		LastSyncTime:    time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := l.nftService.SaveNFT(l.ctx, nft); err != nil {
		return fmt.Errorf("创建NFT记录失败: %v", err)
	}

	log.Printf("⚠️ NFT记录不存在，已创建: %s/%s", contractAddr, tokenID)
	return nil
}

// ---------------- 拍卖事件监听 ----------------
func (l *BlockchainListener) listenAuctionEvents() {
	auctionAddr := l.auctionService.GetContractAddress()
	query := ethereum.FilterQuery{Addresses: []common.Address{auctionAddr}}

	// 提前计算事件签名（只计算一次，提高性能）
	auctionCreatedID := crypto.Keccak256Hash([]byte("AuctionCreated(uint256,address,uint256,uint256)"))
	bidPlacedID := crypto.Keccak256Hash([]byte("NewBid(uint256,address,uint256)"))
	auctionEndedID := crypto.Keccak256Hash([]byte("AuctionEnded(uint256,address,uint256)"))

	// 提前创建Filterer（避免每次循环都创建）
	// 只需要解析，不需要查询, 第二个参数可以传eth client 与区块链交互
	filterer, err := contract.NewNftAuctionFilterer(auctionAddr, nil)
	if err != nil {
		log.Fatalf("❌ 创建Filterer失败: %v", err)
	}

	log.Println("✅ NFT拍卖事件监听器订阅成功，等待事件...")
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
			var eventName string

			// 根据事件ID分流处理
			switch eventID {
			case auctionCreatedID:
				eventName = "AuctionCreated"
				log.Printf("📥 NFT 拍卖事件收到事件，事件名：%s 签名: %s", eventName, eventID.Hex())

				// 只解析AuctionCreated事件
				event, err := filterer.ParseAuctionCreated(vLog)
				if err != nil {
					log.Printf("❌ 解析AuctionCreated失败: %v", err)
					continue
				}
				auctionID = event.AuctionId
				l.handleAuctionCreated(event, vLog)

			case bidPlacedID:
				eventName = "NewBid"
				log.Printf("📥 NFT 拍卖事件收到事件，事件名：%s 签名: %s", eventName, eventID.Hex())

				// 只解析NewBid事件
				event, err := filterer.ParseNewBid(vLog)
				if err != nil {
					log.Printf("❌ 解析NewBid失败: %v", err)
					continue
				}
				auctionID = event.AuctionId
				l.handleNewBid(event, vLog)

			case auctionEndedID:
				eventName = "AuctionEnded"
				log.Printf("📥 NFT 拍卖事件收到事件，事件名：%s 签名: %s", eventName, eventID.Hex())

				// 只解析AuctionEnded事件
				event, err := filterer.ParseAuctionEnded(vLog)
				if err != nil {
					log.Printf("❌ 解析AuctionEnded失败: %v", err)
					continue
				}
				auctionID = event.AuctionId
				l.handleAuctionEnded(event, vLog)

			default:
				log.Printf("⚠️ 未知拍卖事件: %s", eventID.Hex())
				continue
			}

			// 统计计数
			if auctionID != nil {
				log.Printf("🏷️ 拍卖事件: %s, AuctionID=%s", eventName, auctionID.String())

				l.statsLock.Lock()
				l.stats["auctions"]++
				if eventID == bidPlacedID {
					l.stats["bids"]++
				}
				l.statsLock.Unlock()

				// 注意：现在不需要调用 UpdateAuctionFromChain 了！
				// 因为 handleXXX 方法已经用事件数据更新了数据库
			}
		case <-l.ctx.Done():
			log.Println("❌ 拍卖监听器已停止")
			return
		}
	}
}

// 处理拍卖创建事件 - 现在可以直接使用事件参数
func (l *BlockchainListener) handleAuctionCreated(event *contract.NftAuctionAuctionCreated, vLog types.Log) {
	// 直接从事件获取所有参数，不需要再查区块链
	auction := &model.Auction{
		AuctionID:     event.AuctionId.Uint64(),
		NFTContract:   l.auctionService.GetContractAddress().Hex(), // 假设拍卖合约知道对应的NFT合约
		TokenID:       event.TokenId.String(),
		Seller:        event.Seller.Hex(),
		StartingPrice: event.StartPrice.String(),
		HighestBid:    "0",
		HighestBidder: "0x0000000000000000000000000000000000000000",
		StartTime:     uint64(time.Now().Unix()), // 可能需要从区块时间获取更准确
		EndTime:       0,                         // 需要从duration计算，可能需要额外查询
		Ended:         false,
		Status:        "active",
	}

	// 如果有问题，可以记录但不阻塞
	if err := l.auctionService.SaveAuction(l.ctx, auction); err != nil {
		log.Printf("❌ 保存拍卖失败: %v", err)
	} else {
		log.Printf("✅ 拍卖 #%d 已保存到数据库", auction.AuctionID)
	}
}

// 处理新出价事件
func (l *BlockchainListener) handleNewBid(event *contract.NftAuctionNewBid, vLog types.Log) {
	// 1. 保存出价历史
	bidHistory := &model.BidHistory{
		AuctionID:   event.AuctionId.Uint64(),
		Bidder:      event.Bidder.Hex(),
		Amount:      event.Amount.String(),
		TxHash:      vLog.TxHash.Hex(),
		BlockNumber: vLog.BlockNumber,
		BlockTime:   uint64(time.Now().Unix()),
		Status:      "success",
	}

	if err := l.auctionService.SaveBidHistory(l.ctx, bidHistory); err != nil {
		log.Printf("❌ 保存出价历史失败: %v", err)
	}

	// 2. 更新拍卖最高出价
	// 注意：这里最好从数据库获取当前拍卖信息来比较
	auction, err := l.auctionService.GetAuctionByAuctionID(l.ctx, event.AuctionId.Uint64())
	if err != nil {
		log.Printf("❌ 获取拍卖 #%d 信息失败: %v", event.AuctionId.Uint64(), err)
		return
	}

	currentBid, _ := new(big.Int).SetString(auction.HighestBid, 10)
	if event.Amount.Cmp(currentBid) > 0 {
		// 更新为更高的出价
		auction.HighestBid = event.Amount.String()
		auction.HighestBidder = event.Bidder.Hex()
		auction.UpdatedAt = time.Now()

		if err := l.auctionService.SaveAuction(l.ctx, auction); err != nil {
			log.Printf("❌ 更新拍卖出价失败: %v", err)
		} else {
			log.Printf("✅ 拍卖 #%d 最高出价更新为 %s", auction.AuctionID, event.Amount.String())
		}
	}
}

// 处理拍卖结束事件
func (l *BlockchainListener) handleAuctionEnded(event *contract.NftAuctionAuctionEnded, vLog types.Log) {
	// 更新拍卖状态为结束
	auction, err := l.auctionService.GetAuctionByAuctionID(l.ctx, event.AuctionId.Uint64())
	if err != nil {
		log.Printf("❌ 获取拍卖 #%d 信息失败: %v", event.AuctionId.Uint64(), err)
		return
	}

	auction.Ended = true
	auction.Status = "ended"
	auction.HighestBid = event.FinalPrice.String()
	auction.HighestBidder = event.Winner.Hex()
	auction.UpdatedAt = time.Now()

	if err := l.auctionService.SaveAuction(l.ctx, auction); err != nil {
		log.Printf("❌ 更新拍卖结束状态失败: %v", err)
	} else {
		log.Printf("✅ 拍卖 #%d 已结束，赢家: %s", auction.AuctionID, event.Winner.Hex())
	}
}
