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
	// 模拟模式
	if rpcURL == "" || contractAddress == "" {
		log.Println("📡 创建模拟拍卖客户端（演示模式）")
		return &AuctionClient{
			client:   nil,
			contract: nil,
			address:  common.Address{},
			active:   false,
			rpcURL:   "",
		}, nil
	}

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

	if !c.active {
		// 返回模拟数据
		return c.getMockAuctionInfo(auctionID)
	}

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
	if !c.active {
		return big.NewInt(3), nil // 模拟3个拍卖
	}
	return c.contract.NextAuctionId(&bind.CallOpts{Context: ctx})
}

// GetAdmin 获取管理员地址
func (c *AuctionClient) GetAdmin(ctx context.Context) (common.Address, error) {
	if !c.active {
		return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
	}
	return c.contract.Admin(&bind.CallOpts{Context: ctx})
}

// IsTokenAllowed 检查ERC20代币是否被允许
func (c *AuctionClient) IsTokenAllowed(ctx context.Context, tokenAddress common.Address) (bool, error) {
	if !c.active {
		return true, nil
	}
	return c.contract.AllowedERC20Tokens(&bind.CallOpts{Context: ctx}, tokenAddress)
}

// ==================== 交易方法（需要签名）====================
// 注意：这些方法需要私钥签名，这里只提供接口定义

// PlaceBidETH 使用ETH出价
func (c *AuctionClient) PlaceBidETH(ctx context.Context, auctionID *big.Int, amount *big.Int) error {
	if !c.active {
		log.Printf("模拟出价: 拍卖 #%s 出价 %s ETH",
			auctionID.String(), formatWeiToEth(amount))
		return nil
	}
	return fmt.Errorf("出价需要签名交易，请配置私钥")
}

// PlaceBidERC20 使用ERC20出价
func (c *AuctionClient) PlaceBidERC20(ctx context.Context, auctionID *big.Int, amount *big.Int) error {
	if !c.active {
		log.Printf("模拟出价(ERC20): 拍卖 #%s 出价 %s 代币",
			auctionID.String(), amount.String())
		return nil
	}
	return fmt.Errorf("出价需要签名交易，请配置私钥")
}

// EndAuction 结束拍卖
func (c *AuctionClient) EndAuction(ctx context.Context, auctionID *big.Int) error {
	if !c.active {
		log.Printf("模拟结束拍卖: #%s", auctionID.String())
		return nil
	}
	return fmt.Errorf("结束拍卖需要签名交易，请配置私钥")
}

// CreateAuctionETH 创建ETH拍卖
func (c *AuctionClient) CreateAuctionETH(ctx context.Context, duration *big.Int, startPrice *big.Int,
	nftAddress common.Address, tokenID *big.Int) error {
	if !c.active {
		log.Printf("模拟创建ETH拍卖: 时长 %s秒, 起拍价 %s ETH, NFT #%s",
			duration.String(), formatWeiToEth(startPrice), tokenID.String())
		return nil
	}
	return fmt.Errorf("创建拍卖需要签名交易，请配置私钥")
}

// CreateAuctionERC20 创建ERC20拍卖
func (c *AuctionClient) CreateAuctionERC20(ctx context.Context, duration *big.Int, startPrice *big.Int,
	nftAddress common.Address, tokenID *big.Int, erc20Token common.Address) error {
	if !c.active {
		log.Printf("模拟创建ERC20拍卖: 时长 %s秒, 起拍价 %s 代币, NFT #%s",
			duration.String(), startPrice.String(), tokenID.String())
		return nil
	}
	return fmt.Errorf("创建拍卖需要签名交易，请配置私钥")
}

// SetAuctionToken 设置拍卖接受的代币类型
func (c *AuctionClient) SetAuctionToken(ctx context.Context, auctionID *big.Int, tokenAddress common.Address) error {
	if !c.active {
		log.Printf("模拟设置拍卖代币: 拍卖 #%s 接受代币 %s",
			auctionID.String(), tokenAddress.Hex())
		return nil
	}
	return fmt.Errorf("设置代币需要签名交易，请配置私钥")
}

// AllowERC20Token 允许ERC20代币
func (c *AuctionClient) AllowERC20Token(ctx context.Context, tokenAddress common.Address) error {
	if !c.active {
		log.Printf("模拟允许ERC20代币: %s", tokenAddress.Hex())
		return nil
	}
	return fmt.Errorf("允许代币需要签名交易，请配置私钥")
}

// ==================== 辅助方法 ====================

// IsActive 检查客户端是否活跃
func (c *AuctionClient) IsActive() bool {
	return c.active
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

// 获取模拟拍卖数据（保持你原来的风格）
func (c *AuctionClient) GetMockAuctions() ([]struct {
	AuctionID     uint64
	Seller        string
	StartingPrice *big.Int
	HighestBid    *big.Int
	HighestBidder string
}, error) {
	return []struct {
		AuctionID     uint64
		Seller        string
		StartingPrice *big.Int
		HighestBid    *big.Int
		HighestBidder string
	}{
		{
			AuctionID:     1,
			Seller:        "0x1234567890123456789012345678901234567890",
			StartingPrice: big.NewInt(1000000000000000000), // 1 ETH
			HighestBid:    big.NewInt(1200000000000000000), // 1.2 ETH
			HighestBidder: "0x9876543210987654321098765432109876543210",
		},
		{
			AuctionID:     2,
			Seller:        "0x1111111111111111111111111111111111111111",
			StartingPrice: big.NewInt(2500000000000000000), // 2.5 ETH
			HighestBid:    big.NewInt(0),                   // 暂无出价
			HighestBidder: "",
		},
	}, nil
}

// ==================== 私有辅助方法 ====================

func (c *AuctionClient) getMockAuctionInfo(auctionID *big.Int) (
	common.Address, *big.Int, *big.Int, *big.Int, bool, common.Address, *big.Int,
	common.Address, *big.Int, common.Address, *big.Int, *big.Int, error) {

	// 根据拍卖ID返回不同的模拟数据
	switch auctionID.Uint64() {
	case 0:
		return common.HexToAddress("0x1111111111111111111111111111111111111111"),
			big.NewInt(3600), // 1小时
			big.NewInt(1000000000000000000), // 1 ETH
			big.NewInt(time.Now().Unix() - 1800), // 30分钟前开始
			false,
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
			big.NewInt(1500000000000000000), // 1.5 ETH
			common.HexToAddress("0x3333333333333333333333333333333333333333"),
			big.NewInt(1),
			common.Address{}, // ETH拍卖
			big.NewInt(1500000000000000000),
			big.NewInt(1800), // 剩余30分钟
			nil
	case 1:
		return common.HexToAddress("0x4444444444444444444444444444444444444444"),
			big.NewInt(7200), // 2小时
			big.NewInt(5000000000000000000), // 5 ETH
			big.NewInt(time.Now().Unix() - 3600), // 1小时前开始
			false,
			common.Address{}, // 暂无出价者
			big.NewInt(0), // 暂无出价
			common.HexToAddress("0x5555555555555555555555555555555555555555"),
			big.NewInt(2),
			common.Address{}, // ETH拍卖
			big.NewInt(0),
			big.NewInt(3600), // 剩余1小时
			nil
	default:
		return common.Address{}, nil, nil, nil, false, common.Address{}, nil,
			common.Address{}, nil, common.Address{}, nil, nil,
			fmt.Errorf("拍卖不存在")
	}
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
