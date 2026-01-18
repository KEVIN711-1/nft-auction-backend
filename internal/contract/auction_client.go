package contract

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// AuctionClient 拍卖合约客户端
type AuctionClient struct {
	client   *ethclient.Client
	contract *NftAuction
	address  common.Address
	active   bool
	rpcURL   string
}

// NewAuctionClient 创建拍卖客户端
func NewAuctionClient(rpcURL string, contractAddress string) (*AuctionClient, error) {
	log.Printf("正在连接到以太坊节点（拍卖合约）: %s", rpcURL)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接以太坊节点失败: %v", err)
	}

	address := common.HexToAddress(contractAddress)
	contract, err := NewNftAuction(address, client)
	if err != nil {
		return nil, fmt.Errorf("初始化拍卖合约失败: %v", err)
	}

	// 测试连接
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("测试网络连接失败: %v", err)
	}

	log.Printf("✅ 拍卖合约连接成功，网络ID: %v", networkID)
	log.Printf("✅ 拍卖合约地址: %s", address.Hex())

	return &AuctionClient{
		client:   client,
		contract: contract,
		address:  address,
		active:   true,
		rpcURL:   rpcURL,
	}, nil
}

// ==================== 查询方法（不需要签名）====================

// GetAuctionInfo 获取拍卖详细信息
func (c *AuctionClient) GetAuctionInfo(ctx context.Context, auctionID *big.Int) (
	common.Address, *big.Int, *big.Int, *big.Int, bool, common.Address, *big.Int,
	common.Address, *big.Int, common.Address, *big.Int, *big.Int, error) {

	// 调用拍卖合约的 auctions 映射
	auction, err := c.contract.Auctions(&bind.CallOpts{Context: ctx}, auctionID)
	if err != nil {
		return common.Address{}, nil, nil, nil, false, common.Address{}, nil,
			common.Address{}, nil, common.Address{}, nil, nil, err
	}

	// 根据实际结构体字段获取数据
	tokenAddress := common.Address{}
	bidTokenAmount := big.NewInt(0)

	// 检查实际字段名（根据你的合约）
	// 如果合约有 UseERC20 和 Erc20Token 字段
	if auction.UseERC20 {
		tokenAddress = auction.Erc20Token
		// 这里需要根据实际情况获取 bidTokenAmount
		// 可能需要从其他地方获取或使用默认值
	}

	// 计算剩余时间
	timeRemaining := c.calculateTimeRemaining(auction)

	return auction.Seller,
		auction.Duration,
		auction.StartPrice,
		auction.StartTime,
		auction.Ended,
		auction.HighestBidder,
		auction.HighestBid,
		auction.NftContract,
		auction.TokenId,
		tokenAddress, // 使用调整后的值
		bidTokenAmount, // 使用调整后的值
		timeRemaining,
		nil
}

// GetAuctionCount 获取拍卖总数
func (c *AuctionClient) GetAuctionCount(ctx context.Context) (*big.Int, error) {
	return c.contract.NextAuctionId(&bind.CallOpts{Context: ctx})
}

// GetAdmin 获取管理员地址
func (c *AuctionClient) GetAdmin(ctx context.Context) (common.Address, error) {
	return c.contract.Admin(&bind.CallOpts{Context: ctx})
}

// IsTokenAllowed 检查ERC20代币是否被允许
func (c *AuctionClient) IsTokenAllowed(ctx context.Context, tokenAddress common.Address) (bool, error) {
	return c.contract.AllowedERC20Tokens(&bind.CallOpts{Context: ctx}, tokenAddress)
}

// GetContractAddress 获取合约地址
func (c *AuctionClient) GetContractAddress() common.Address {
	return c.address
}

// GetLatestBlockNumber 获取最新区块号
func (c *AuctionClient) GetLatestBlockNumber() (uint64, error) {
	if !c.active {
		return 12345678, nil
	}
	header, err := c.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, fmt.Errorf("获取区块信息失败: %v", err)
	}
	return header.Number.Uint64(), nil
}

// 修改 calculateTimeRemaining 函数，使用实际的合约结构体
func (c *AuctionClient) calculateTimeRemaining(auction struct {
	Seller        common.Address
	Duration      *big.Int
	StartPrice    *big.Int
	StartTime     *big.Int
	Ended         bool
	HighestBidder common.Address
	HighestBid    *big.Int
	NftContract   common.Address
	TokenId       *big.Int
	UseERC20      bool
	Erc20Token    common.Address
}) *big.Int {
	if auction.Ended || auction.StartTime == nil || auction.Duration == nil {
		return big.NewInt(0)
	}

	startTime := auction.StartTime.Uint64()
	duration := auction.Duration.Uint64()
	currentTime := uint64(time.Now().Unix())

	if startTime+duration <= currentTime {
		return big.NewInt(0)
	}

	return big.NewInt(int64(startTime + duration - currentTime))
}

// 辅助函数：格式化wei为ETH
func formatWeiToEth(wei *big.Int) string {
	if wei == nil {
		return "0"
	}
	eth := new(big.Float).SetInt(wei)
	eth = eth.Quo(eth, big.NewFloat(1e18))
	return eth.Text('f', 4)
}
