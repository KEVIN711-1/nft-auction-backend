package contract

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 你的智能合约代码 (.sol)
//         │
//         ▼ 编译、部署
// ┌─────────────────────┐
// │  以太坊区块链网络   │
// │  KevinNFT 合约      │
// └─────────────────────┘
//         │
//         ▼ abigen工具生成
// ┌─────────────────────┐
// │   kevinnft.go       │  ← 自动生成的Go绑定文件
// │   - 合约的Go包装器  │
// │   - 包含所有ABI方法 │
// │   - 类型安全调用    │
// └─────────────────────┘
//         │
//         ▼ 被调用
// ┌─────────────────────┐
// │   nft_client.go     │  ← 你写的客户端逻辑
// │   - 连接以太坊节点 │
// │   - 管理连接状态   │
// │   - 调用kevinnft.go│
// │   - 错误处理       │
// └─────────────────────┘
//         │
//         ▼ 实现接口
// ┌─────────────────────┐
// │   contract.go       │  ← 接口定义
// │   - 定义方法签名    │
// │   - 抽象层         │
// └─────────────────────┘

// NFTClient 现在使用你的 KevinNFT 合约
type NFTClient struct {
	client   *ethclient.Client
	contract *KevinNFT
	address  common.Address
}

// NewNFTClient 创建新的 NFT 客户端
func NewNFTClient(rpcURL string, contractAddress string) (*NFTClient, error) {
	log.Printf("正在连接到以太坊节点（NFT合约）: %s", rpcURL)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	address := common.HexToAddress(contractAddress)
	contract, err := NewKevinNFT(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract: %v", err)
	}

	// 测试连接
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("测试网络连接失败: %v", err)
	}

	log.Printf("✅ NFT合约连接成功，网络ID: %v", networkID)
	log.Printf("✅ NFT合约地址: %s", address.Hex())

	return &NFTClient{
		client:   client,
		contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress 获取合约地址 - 🔥 新增方法
func (c *NFTClient) GetContractAddress() common.Address {
	return c.address
}

// GetName 获取合约名称
func (c *NFTClient) GetName(ctx context.Context) (string, error) {
	return c.contract.Name(&bind.CallOpts{Context: ctx})
}

// GetSymbol 获取合约符号
func (c *NFTClient) GetSymbol(ctx context.Context) (string, error) {
	return c.contract.Symbol(&bind.CallOpts{Context: ctx})
}

// GetOwner 获取 NFT 所有者
func (c *NFTClient) GetOwner(ctx context.Context, tokenID *big.Int) (common.Address, error) {
	return c.contract.OwnerOf(&bind.CallOpts{Context: ctx}, tokenID)
}

// GetTokenURI 获取 token URI
func (c *NFTClient) GetTokenURI(ctx context.Context, tokenID *big.Int) (string, error) {
	return c.contract.TokenURI(&bind.CallOpts{Context: ctx}, tokenID)
}

// CheckOwner 检查用户是否是 NFT 所有者
func (c *NFTClient) CheckOwner(ctx context.Context, tokenID *big.Int, address string) (bool, error) {
	owner, err := c.GetOwner(ctx, tokenID)
	if err != nil {
		return false, err
	}

	checkAddr := common.HexToAddress(address)
	return owner.Hex() == checkAddr.Hex(), nil
}

// GetTotalSupply 获取总供应量（需要合约支持）
func (c *NFTClient) GetTotalSupply(ctx context.Context) (*big.Int, error) {
	maxTokenID := big.NewInt(10) // 设置一个合理的上限

	foundCount := big.NewInt(0)

	for i := int64(1); i < maxTokenID.Int64(); i++ {
		tokenID := big.NewInt(i)

		// 检查NFT是否存在
		exists, _ := c.CheckIfMinted(ctx, tokenID)
		if exists {
			foundCount.Add(foundCount, big.NewInt(1))
		} else {
			break
		}
	}

	return foundCount, nil
}

// GetBalanceOf 获取地址拥有的 NFT 数量
func (c *NFTClient) GetBalanceOf(ctx context.Context, address common.Address) (*big.Int, error) {
	return c.contract.BalanceOf(&bind.CallOpts{Context: ctx}, address)
}

// CheckIfMinted 检查 NFT 是否已被铸造
func (c *NFTClient) CheckIfMinted(ctx context.Context, tokenID *big.Int) (bool, error) {
	_, err := c.contract.OwnerOf(&bind.CallOpts{Context: ctx}, tokenID)
	if err != nil {
		// 如果 token 不存在，会返回错误
		if err.Error() == "execution reverted" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
