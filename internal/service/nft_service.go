package service

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"nft-auction-backend/internal/contract"
	"nft-auction-backend/internal/model"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type NFTService struct {
	DB     *gorm.DB
	client contract.NFTContract
}

func NewNFTService(db *gorm.DB, client contract.NFTContract) *NFTService {
	return &NFTService{
		DB:     db,
		client: client,
	}
}

// nft_service.go - 添加这个方法
func (s *NFTService) GetContractAddress() common.Address {
	return s.client.GetContractAddress()
}

// GetNFT 从数据库获取 NFT
func (s *NFTService) GetNFT(contractAddr, tokenID string) (*model.NFTInfo, error) {
	var nft model.NFTInfo
	result := s.DB.Where("contract_address = ? AND token_id = ?", contractAddr, tokenID).First(&nft)
	if result.Error != nil {
		return nil, result.Error
	}
	return &nft, nil
}

// SaveNFT 保存或更新 NFT 到数据库
func (s *NFTService) SaveNFT(ctx context.Context, nft *model.NFTInfo) error {
	if nft == nil {
		return nil
	}

	var existing model.NFTInfo
	result := s.DB.WithContext(ctx).
		Where("contract_address = ? AND token_id = ?", nft.ContractAddress, nft.TokenID).
		First(&existing)
	now := time.Now()

	if result.Error != nil {
		nft.CreatedAt = now
		nft.UpdatedAt = now
		if err := s.DB.WithContext(ctx).Create(nft).Error; err != nil {
			log.Printf("❌ 保存 NFT token id = %s 失败: %v", nft.TokenID, err)
			return err
		}
		log.Printf("新增 NFT %s", nft.TokenID)
		return nil
	}

	// 更新现有记录
	existing.Owner = nft.Owner
	existing.TotalSupply = nft.TotalSupply
	existing.Blockchain = nft.Blockchain
	existing.LastSyncTime = now
	existing.ContractName = nft.ContractName
	existing.ContractSymbol = nft.ContractSymbol
	existing.IsMinted = nft.IsMinted
	existing.UpdatedAt = now

	if err := s.DB.WithContext(ctx).Save(&existing).Error; err != nil {
		log.Printf("❌ 更新 NFT %s 失败: %v", existing.TokenID, err)
		return err
	}
	log.Printf(" 更新 NFT token id = %s", existing.TokenID)
	return nil
}

// GetOwner 从区块链获取NFT拥有者
func (s *NFTService) GetOwner(ctx context.Context, tokenID string) (string, error) {
	tokenIDBig, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return "", fmt.Errorf("invalid token ID format: %s", tokenID)
	}

	ownerAddr, err := s.client.GetOwner(ctx, tokenIDBig)
	if err != nil {
		return "", fmt.Errorf("failed to get owner from blockchain: %v", err)
	}

	return ownerAddr.Hex(), nil
}

// ValidateOwnership 验证指定地址是否是NFT所有者
// 优化为从数据库判断
func (s *NFTService) ValidateOwnership(ctx context.Context, tokenID, address string) (bool, error) {
	owner, err := s.GetOwner(ctx, tokenID)
	if err != nil {
		return false, err
	}
	return owner == address, nil
}

// SyncAllNFTs 同步所有NFT
func (s *NFTService) SyncAllNFTs(ctx context.Context) error {
	total, err := s.client.GetTotalSupply(ctx)
	if err != nil {
		return fmt.Errorf("获取总供应量失败: %v", err)
	}

	// 获取合约信息
	contractAddr := s.client.GetContractAddress().Hex()
	contractName, _ := s.client.GetName(ctx)
	contractSymbol, _ := s.client.GetSymbol(ctx)

	successCount := 0
	for i := int64(1); i <= total.Int64(); i++ {
		// 检查ctx是否已取消
		select {
		case <-ctx.Done():
			return fmt.Errorf("同步被取消: %v", ctx.Err())
		default:
			// 继续执行
		}

		tokenID := fmt.Sprintf("%d", i)
		tokenIDBig := big.NewInt(i)

		ownerAddr, err := s.client.GetOwner(ctx, tokenIDBig)
		if err != nil {
			log.Printf("❌ 获取 NFT %s 所有者失败: %v", tokenID, err)
			continue
		}

		nft := &model.NFTInfo{
			ContractAddress: contractAddr,
			TokenID:         tokenID,
			Owner:           ownerAddr.Hex(),
			Name:            fmt.Sprintf("NFT #%s", tokenID),
			TotalSupply:     total.String(),
			Blockchain:      "sepolia",
			ContractName:    contractName,
			ContractSymbol:  contractSymbol,
			IsMinted:        true,
			LastSyncTime:    time.Now(),
		}

		// 注意：这里需要传递ctx给SaveNFT
		if s.SaveNFT(ctx, nft) == nil {
			successCount++
		}
	}

	log.Printf(" NFT全量同步完成，总数: %d", total.Int64())
	return nil
}

// UpdateNFTFromChain 从链上更新单个NFT信息（事件监听器调用）
func (s *NFTService) UpdateNFTFromChain(ctx context.Context, tokenID string) error {
	// 注意：不再创建新的context，使用传入的ctx

	// 将 tokenID 转换为 big.Int
	tokenIDBig, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return fmt.Errorf("invalid token ID format: %s", tokenID)
	}

	// 从区块链获取所有者
	ownerAddr, err := s.client.GetOwner(ctx, tokenIDBig)
	if err != nil {
		return fmt.Errorf("failed to get NFT owner from blockchain: %v", err)
	}

	// 获取合约信息
	contractAddr := s.client.GetContractAddress().Hex()
	contractName, _ := s.client.GetName(ctx)
	contractSymbol, _ := s.client.GetSymbol(ctx)
	contractUrl, _ := s.client.GetTokenURI(ctx, tokenIDBig)

	// 获取总供应量
	var totalSupply string
	if total, err := s.client.GetTotalSupply(ctx); err == nil {
		totalSupply = total.String()
	}

	// 构建NFT信息
	nft := &model.NFTInfo{
		ContractAddress: contractAddr,
		TokenID:         tokenID,
		Owner:           ownerAddr.Hex(),
		Name:            fmt.Sprintf("NFT #%s", tokenID),
		Uri:             contractUrl, // 修正：直接使用获取到的URI
		TotalSupply:     totalSupply,
		Blockchain:      "sepolia",
		ContractName:    contractName,
		ContractSymbol:  contractSymbol,
		IsMinted:        true,
		LastSyncTime:    time.Now(),
	}

	// 保存到数据库，传递ctx
	err = s.SaveNFT(ctx, nft)
	if err != nil {
		return fmt.Errorf("保存NFT到数据库失败: %s/%s", contractAddr, tokenID)
	}

	log.Printf(" NFT已更新: %s/%s, Owner=%s", contractAddr, tokenID, ownerAddr.Hex())
	return nil
}
