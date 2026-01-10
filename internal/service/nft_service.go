// internal/service/nft_service.go
package service

import (
	"fmt"
	"log"
	"math/big"

	"nft-auction-backend/internal/contract"
	"nft-auction-backend/internal/model"

	"gorm.io/gorm"
)

type NFTService struct {
	db        *gorm.DB
	nftClient *contract.NFTClient
}

func NewNFTService(db *gorm.DB, nftClient *contract.NFTClient) *NFTService {
	return &NFTService{
		db:        db,
		nftClient: nftClient,
	}
}

// SyncNFTInfo 同步NFT合约信息到数据库
func (s *NFTService) SyncNFTInfo() error {
	log.Println("🔄 开始同步NFT合约信息...")

	// 1. 获取NFT名称
	name, err := s.nftClient.GetName()
	if err != nil {
		return fmt.Errorf("获取NFT名称失败: %v", err)
	}

	// 2. 获取NFT符号
	symbol, err := s.nftClient.GetSymbol()
	if err != nil {
		return fmt.Errorf("获取NFT符号失败: %v", err)
	}

	// 3. 获取总供应量
	totalSupply, err := s.nftClient.GetTotalSupply()
	if err != nil {
		// 有些合约可能没有totalSupply方法，设为0
		log.Printf("⚠️  获取总供应量失败: %v，使用默认值0", err)
		totalSupply = big.NewInt(0)
	}

	// 4. 保存到数据库
	nftInfo := model.NFTInfo{
		ContractAddress: s.nftClient.GetContractAddress(), // 需要添加这个方法
		Name:            name,
		Symbol:          symbol,
		TotalSupply:     totalSupply.String(),
		Owner:           "", // 后续可以添加获取合约所有者的方法
		Blockchain:      "sepolia",
	}

	// 检查是否已存在
	var existing model.NFTInfo
	result := s.db.Where("contract_address = ?", nftInfo.ContractAddress).First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新记录
		if err := s.db.Create(&nftInfo).Error; err != nil {
			return fmt.Errorf("创建NFT信息失败: %v", err)
		}
		log.Printf("✅ 创建NFT信息: %s (%s)", name, symbol)
	} else if result.Error == nil {
		// 更新现有记录
		if err := s.db.Model(&existing).Updates(&nftInfo).Error; err != nil {
			return fmt.Errorf("更新NFT信息失败: %v", err)
		}
		log.Printf("🔄 更新NFT信息: %s (%s)", name, symbol)
	} else {
		return fmt.Errorf("查询NFT信息失败: %v", result.Error)
	}

	log.Println("✅ NFT合约信息同步完成")
	return nil
}

// GetNFTInfo 获取NFT信息
func (s *NFTService) GetNFTInfo() (*model.NFTInfo, error) {
	var nftInfo model.NFTInfo

	// 先尝试从数据库获取
	result := s.db.First(&nftInfo)
	if result.Error != nil {
		return nil, fmt.Errorf("获取NFT信息失败: %v", result.Error)
	}

	return &nftInfo, nil
}

// GetNFTInfoByAddress 根据合约地址获取NFT信息
func (s *NFTService) GetNFTInfoByAddress(contractAddress string) (*model.NFTInfo, error) {
	var nftInfo model.NFTInfo
	result := s.db.Where("contract_address = ?", contractAddress).First(&nftInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &nftInfo, nil
}
