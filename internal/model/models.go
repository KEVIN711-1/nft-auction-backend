package model

import (
	"time"
)

// Auction 拍卖表
type Auction struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement"`
	AuctionID     uint64 `gorm:"uniqueIndex"` // 合约中的拍卖ID
	NFTContract   string `gorm:"size:42"`
	TokenID       string
	Seller        string `gorm:"size:42"`
	StartingPrice string
	HighestBid    string
	HighestBidder string `gorm:"size:42"`
	StartTime     uint64
	EndTime       uint64
	Ended         bool
	TxHash        string    `gorm:"size:66"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// 在 internal/model/models.go 中添加
type NFTInfo struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"`
	ContractAddress string    `gorm:"uniqueIndex;size:42;comment:NFT合约地址"`
	Name            string    `gorm:"size:255;comment:NFT名称"`
	Symbol          string    `gorm:"size:50;comment:NFT符号"`
	TotalSupply     string    `gorm:"type:varchar(100);comment:总供应量"`
	Owner           string    `gorm:"size:42;comment:合约所有者"`
	Blockchain      string    `gorm:"size:20;default:'sepolia';comment:区块链网络"` // 确保有这个字段
	LastSyncTime    time.Time `gorm:"comment:最后同步时间"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (Auction) TableName() string {
	return "auctions"
}
