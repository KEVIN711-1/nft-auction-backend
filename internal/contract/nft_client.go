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

// NFTClient 现在使用你的 KevinNFT 合约
type NFTClient struct {
	client   *ethclient.Client
	contract *KevinNFT
	address  common.Address
}

// NewNFTClient 创建新的 NFT 客户端
func NewNFTClient(rpcURL string, contractAddress string) (*NFTClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	address := common.HexToAddress(contractAddress)
	contract, err := NewKevinNFT(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract: %v", err)
	}

	return &NFTClient{
		client:   client,
		contract: contract,
		address:  address,
	}, nil
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

// TransferFrom 转移 NFT（需要已授权）
func (c *NFTClient) TransferFrom(ctx context.Context, from, to common.Address, tokenID *big.Int) error {
	// 这里需要私钥签名交易
	// 实际实现需要配置私钥
	log.Printf("Transfer NFT %s from %s to %s", tokenID.String(), from.Hex(), to.Hex())
	return nil
}

// GetTotalSupply 获取总供应量（需要合约支持）
func (c *NFTClient) GetTotalSupply(ctx context.Context) (*big.Int, error) {
	// 注意：你的合约目前没有 totalSupply 函数
	// 如果需要，可以在合约中添加
	return big.NewInt(0), nil
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
