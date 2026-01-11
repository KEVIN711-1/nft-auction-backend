package service

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"nft-auction-backend/internal/contract"

	"github.com/ethereum/go-ethereum/common"
)

type NFTService struct {
	client contract.NFTContract
}

func NewNFTService(client contract.NFTContract) *NFTService {
	return &NFTService{
		client: client,
	}
}

// GetNFTInfo 获取 NFT 信息
func (s *NFTService) GetNFTInfo(ctx context.Context, tokenID string) (*NFTInfo, error) {
	// 转换 tokenID
	id, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return nil, fmt.Errorf("invalid token ID: %s", tokenID)
	}

	// 检查 NFT 是否存在
	minted, err := s.client.CheckIfMinted(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check NFT: %v", err)
	}

	if !minted {
		return nil, fmt.Errorf("NFT %s does not exist", tokenID)
	}

	// 获取 NFT 所有者
	owner, err := s.client.GetOwner(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner: %v", err)
	}

	// 获取 Token URI
	tokenURI, err := s.client.GetTokenURI(ctx, id)
	if err != nil {
		log.Printf("Warning: failed to get token URI: %v", err)
		// 继续执行，tokenURI 可能为空
	}

	// 获取合约信息
	contractName, err := s.client.GetName(ctx)
	if err != nil {
		log.Printf("Warning: failed to get contract name: %v", err)
	}

	contractSymbol, err := s.client.GetSymbol(ctx)
	if err != nil {
		log.Printf("Warning: failed to get contract symbol: %v", err)
	}

	return &NFTInfo{
		TokenID:        tokenID,
		Owner:          owner.Hex(),
		TokenURI:       tokenURI,
		ContractName:   contractName,
		ContractSymbol: contractSymbol,
		IsMinted:       true,
	}, nil
}

// ValidateOwnership 验证 NFT 所有权
func (s *NFTService) ValidateOwnership(ctx context.Context, tokenID, address string) (bool, error) {
	// 清理地址格式
	address = strings.ToLower(strings.TrimSpace(address))
	if !common.IsHexAddress(address) {
		return false, fmt.Errorf("invalid Ethereum address: %s", address)
	}

	// 转换 tokenID
	id, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return false, fmt.Errorf("invalid token ID: %s", tokenID)
	}

	return s.client.CheckOwner(ctx, id, address)
}

// GetOwner 获取 NFT 所有者
func (s *NFTService) GetOwner(ctx context.Context, tokenID string) (string, error) {
	id, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return "", fmt.Errorf("invalid token ID: %s", tokenID)
	}

	owner, err := s.client.GetOwner(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get owner: %v", err)
	}

	return owner.Hex(), nil
}

// TransferNFT 转移 NFT（简化版，实际需要私钥签名）
func (s *NFTService) TransferNFT(ctx context.Context, from, to, tokenID string) error {
	// 验证地址
	if !common.IsHexAddress(from) || !common.IsHexAddress(to) {
		return fmt.Errorf("invalid Ethereum address")
	}

	// 转换 tokenID
	id, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return fmt.Errorf("invalid token ID: %s", tokenID)
	}

	// 验证发送者确实是所有者
	isOwner, err := s.ValidateOwnership(ctx, tokenID, from)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("address %s is not the owner of NFT %s", from, tokenID)
	}

	fromAddr := common.HexToAddress(from)
	toAddr := common.HexToAddress(to)

	return s.client.TransferFrom(ctx, fromAddr, toAddr, id)
}

// NFTInfo 结构体
type NFTInfo struct {
	TokenID        string `json:"token_id"`
	Owner          string `json:"owner"`
	TokenURI       string `json:"token_uri"`
	ContractName   string `json:"contract_name"`
	ContractSymbol string `json:"contract_symbol"`
	IsMinted       bool   `json:"is_minted"`
}

// // internal/service/nft_service.go
// package service

// import (
// 	"fmt"
// 	"log"
// 	"math/big"

// 	"nft-auction-backend/internal/contract"
// 	"nft-auction-backend/internal/model"

// 	"gorm.io/gorm"
// )

// type NFTService struct {
// 	db        *gorm.DB
// 	nftClient *contract.NFTClient
// }

// func NewNFTService(db *gorm.DB, nftClient *contract.NFTClient) *NFTService {
// 	return &NFTService{
// 		db:        db,
// 		nftClient: nftClient,
// 	}
// }

// // SyncNFTInfo 同步NFT合约信息到数据库
// func (s *NFTService) SyncNFTInfo() error {
// 	log.Println("🔄 开始同步NFT合约信息...")

// 	// 1. 获取NFT名称
// 	name, err := s.nftClient.GetName()
// 	if err != nil {
// 		return fmt.Errorf("获取NFT名称失败: %v", err)
// 	}

// 	// 2. 获取NFT符号
// 	symbol, err := s.nftClient.GetSymbol()
// 	if err != nil {
// 		return fmt.Errorf("获取NFT符号失败: %v", err)
// 	}

// 	// 3. 获取总供应量
// 	totalSupply, err := s.nftClient.GetTotalSupply()
// 	if err != nil {
// 		// 有些合约可能没有totalSupply方法，设为0
// 		log.Printf("⚠️  获取总供应量失败: %v，使用默认值0", err)
// 		totalSupply = big.NewInt(0)
// 	}

// 	// 4. 保存到数据库
// 	nftInfo := model.NFTInfo{
// 		ContractAddress: s.nftClient.GetContractAddress(), // 需要添加这个方法
// 		Name:            name,
// 		Symbol:          symbol,
// 		TotalSupply:     totalSupply.String(),
// 		Owner:           "", // 后续可以添加获取合约所有者的方法
// 		Blockchain:      "sepolia",
// 	}

// 	// 检查是否已存在
// 	var existing model.NFTInfo
// 	result := s.db.Where("contract_address = ?", nftInfo.ContractAddress).First(&existing)

// 	if result.Error == gorm.ErrRecordNotFound {
// 		// 创建新记录
// 		if err := s.db.Create(&nftInfo).Error; err != nil {
// 			return fmt.Errorf("创建NFT信息失败: %v", err)
// 		}
// 		log.Printf("✅ 创建NFT信息: %s (%s)", name, symbol)
// 	} else if result.Error == nil {
// 		// 更新现有记录
// 		if err := s.db.Model(&existing).Updates(&nftInfo).Error; err != nil {
// 			return fmt.Errorf("更新NFT信息失败: %v", err)
// 		}
// 		log.Printf("🔄 更新NFT信息: %s (%s)", name, symbol)
// 	} else {
// 		return fmt.Errorf("查询NFT信息失败: %v", result.Error)
// 	}

// 	log.Println("✅ NFT合约信息同步完成")
// 	return nil
// }

// // GetNFTInfo 获取NFT信息
// func (s *NFTService) GetNFTInfo() (*model.NFTInfo, error) {
// 	var nftInfo model.NFTInfo

// 	// 先尝试从数据库获取
// 	result := s.db.First(&nftInfo)
// 	if result.Error != nil {
// 		return nil, fmt.Errorf("获取NFT信息失败: %v", result.Error)
// 	}

// 	return &nftInfo, nil
// }

// // GetNFTInfoByAddress 根据合约地址获取NFT信息
// func (s *NFTService) GetNFTInfoByAddress(contractAddress string) (*model.NFTInfo, error) {
// 	var nftInfo model.NFTInfo
// 	result := s.db.Where("contract_address = ?", contractAddress).First(&nftInfo)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &nftInfo, nil
// }
