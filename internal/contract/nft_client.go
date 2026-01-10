package contract

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type NFTClient struct {
	client *ethclient.Client
	abi    abi.ABI
	addr   common.Address
}

func NewNFTClient(rpcURL, contractAddress string) (*NFTClient, error) {
	if rpcURL == "" || contractAddress == "" {
		return nil, fmt.Errorf("RPC URL或合约地址为空")
	}

	log.Printf("正在创建NFT客户端: %s", contractAddress)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接以太坊节点失败: %v", err)
	}

	// 简化的ERC721 ABI（包含常用的方法）
	// 替换为MyNFT.sol 的ABI
	const erc721ABI = `[
		{
			"constant": true,
			"inputs": [],
			"name": "name",
			"outputs": [{"name": "", "type": "string"}],
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "symbol",
			"outputs": [{"name": "", "type": "string"}],
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [],
			"name": "totalSupply",
			"outputs": [{"name": "", "type": "uint256"}],
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [{"name": "tokenId", "type": "uint256"}],
			"name": "ownerOf",
			"outputs": [{"name": "", "type": "address"}],
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [{"name": "owner", "type": "address"}],
			"name": "balanceOf",
			"outputs": [{"name": "", "type": "uint256"}],
			"type": "function"
		}
	]`

	parsedABI, err := abi.JSON(strings.NewReader(erc721ABI))
	if err != nil {
		return nil, fmt.Errorf("解析ABI失败: %v", err)
	}

	return &NFTClient{
		client: client,
		abi:    parsedABI,
		addr:   common.HexToAddress(contractAddress),
	}, nil
}

// GetName 获取NFT名称
func (c *NFTClient) GetName() (string, error) {
	data, err := c.abi.Pack("name")
	if err != nil {
		return "", fmt.Errorf("打包name函数失败: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &c.addr,
		Data: data,
	}

	result, err := c.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", fmt.Errorf("调用name函数失败: %v", err)
	}

	var name string
	err = c.abi.UnpackIntoInterface(&name, "name", result)
	return name, err
}

// GetSymbol 获取NFT符号
func (c *NFTClient) GetSymbol() (string, error) {
	data, err := c.abi.Pack("symbol")
	if err != nil {
		return "", fmt.Errorf("打包symbol函数失败: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &c.addr,
		Data: data,
	}

	result, err := c.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", fmt.Errorf("调用symbol函数失败: %v", err)
	}

	var symbol string
	err = c.abi.UnpackIntoInterface(&symbol, "symbol", result)
	return symbol, err
}

// GetTotalSupply 获取总供应量
func (c *NFTClient) GetTotalSupply() (*big.Int, error) {
	data, err := c.abi.Pack("totalSupply")
	if err != nil {
		return nil, fmt.Errorf("打包totalSupply函数失败: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &c.addr,
		Data: data,
	}

	result, err := c.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("调用totalSupply函数失败: %v", err)
	}

	var totalSupply *big.Int
	err = c.abi.UnpackIntoInterface(&totalSupply, "totalSupply", result)
	return totalSupply, err
}

// GetOwnerOf 获取NFT所有者
func (c *NFTClient) GetOwnerOf(tokenId *big.Int) (common.Address, error) {
	data, err := c.abi.Pack("ownerOf", tokenId)
	if err != nil {
		return common.Address{}, fmt.Errorf("打包ownerOf函数失败: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &c.addr,
		Data: data,
	}

	result, err := c.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return common.Address{}, fmt.Errorf("调用ownerOf函数失败: %v", err)
	}

	var owner common.Address
	err = c.abi.UnpackIntoInterface(&owner, "ownerOf", result)
	return owner, err
}

// GetBalanceOf 获取地址拥有的NFT数量
func (c *NFTClient) GetBalanceOf(owner common.Address) (*big.Int, error) {
	data, err := c.abi.Pack("balanceOf", owner)
	if err != nil {
		return nil, fmt.Errorf("打包balanceOf函数失败: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &c.addr,
		Data: data,
	}

	result, err := c.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("调用balanceOf函数失败: %v", err)
	}

	var balance *big.Int
	err = c.abi.UnpackIntoInterface(&balance, "balanceOf", result)
	return balance, err
}

// 在 internal/contract/nft_client.go 中添加
// GetContractAddress 获取合约地址
func (c *NFTClient) GetContractAddress() string {
	return c.addr.Hex()
}
