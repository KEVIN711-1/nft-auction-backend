package contract

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

// ContractClient 区块链客户端
type ContractClient struct {
	client *ethclient.Client
	rpcURL string
	active bool
}

// NewContractClient 创建区块链客户端
func NewContractClient(rpcURL string) (*ContractClient, error) {
	// 如果RPC URL为空，返回模拟客户端
	if rpcURL == "" {
		log.Println("📡 创建模拟区块链客户端（演示模式）")
		return &ContractClient{
			client: nil,
			rpcURL: "",
			active: false,
		}, nil
	}

	log.Printf("正在连接到以太坊节点: %s", rpcURL)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接以太坊节点失败: %v", err)
	}

	// 测试连接
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("测试网络连接失败: %v", err)
	}

	log.Printf("✅ 连接成功，网络ID: %v", networkID)
	return &ContractClient{
		client: client,
		rpcURL: rpcURL,
		active: true,
	}, nil
}

// IsActive 检查客户端是否活跃（连接到真实节点）
func (c *ContractClient) IsActive() bool {
	return c.active
}

// GetLatestBlockNumber 获取最新区块号
func (c *ContractClient) GetLatestBlockNumber() (uint64, error) {
	if c.client == nil {
		// 返回模拟区块号
		log.Println("📡 使用模拟区块号")
		return 12345678, nil
	}

	header, err := c.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, fmt.Errorf("获取区块信息失败: %v", err)
	}
	return header.Number.Uint64(), nil
}

// GetMockAuctions 模拟获取拍卖数据
func (c *ContractClient) GetMockAuctions() ([]struct {
	AuctionID     uint64
	Seller        string
	StartingPrice *big.Int
	HighestBid    *big.Int
	HighestBidder string
}, error) {
	// 模拟数据
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
			HighestBid:    big.NewInt(0),
			HighestBidder: "",
		},
		{
			AuctionID:     3,
			Seller:        "0x2222222222222222222222222222222222222222",
			StartingPrice: big.NewInt(500000000000000000), // 0.5 ETH
			HighestBid:    big.NewInt(800000000000000000), // 0.8 ETH
			HighestBidder: "0x3333333333333333333333333333333333333333",
		},
	}, nil
}
