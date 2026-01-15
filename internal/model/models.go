package model

import (
	"time"

	"gorm.io/gorm"
)

// mapstructure 标签
// 用途：定义配置映射关系，告诉 mapstructure 库如何将配置文件映射到结构体。

// gorm 标签
//  用途：定义数据库表结构，告诉 GORM ORM 如何创建和管理数据库表。

// Auction 拍卖表
type Auction struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement"` // primaryKey: 主键字段，唯一标识记录 | autoIncrement: 自动递增，每次插入自动增加
	AuctionID     uint64 `gorm:"uniqueIndex"`              // uniqueIndex: 创建唯一索引，防止重复值
	NFTContract   string `gorm:"size:42"`                  // size:42: 字符串最大长度42个字符（以太坊地址长度）
	TokenID       string
	Seller        string `gorm:"size:42"` // size:42: 字符串最大长度42个字符（以太坊地址长度）
	StartingPrice string
	HighestBid    string
	HighestBidder string `gorm:"size:42"` // size:42: 字符串最大长度42个字符（以太坊地址长度）
	StartTime     uint64
	EndTime       uint64
	Ended         bool
	TxHash        string    `gorm:"size:66"`        // size:66: 字符串最大长度66个字符（以太坊交易哈希长度）
	CreatedAt     time.Time `gorm:"autoCreateTime"` // autoCreateTime: 自动设置创建时间，记录插入时自动填充
	UpdatedAt     time.Time `gorm:"autoUpdateTime"` // autoUpdateTime: 自动更新时间，记录修改时自动更新
}

// NFTInfo NFT合约信息表
type NFTInfo struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"`                // primaryKey: 主键字段 | autoIncrement: 自动递增
	ContractAddress string    `gorm:"uniqueIndex;size:42;comment:NFT合约地址"`     // uniqueIndex: 创建唯一索引 | size:42: 地址长度限制 | comment: 数据库字段注释
	Name            string    `gorm:"size:255;comment:NFT名称"`                  // size:255: 名称最大长度 | comment: 数据库字段注释
	Symbol          string    `gorm:"size:50;comment:NFT符号"`                   // size:50: 符号最大长度 | comment: 数据库字段注释
	TotalSupply     string    `gorm:"type:varchar(100);comment:总供应量"`          // type:varchar(100): 指定数据库类型为varchar(100) | comment: 数据库字段注释
	Owner           string    `gorm:"size:42;comment:合约所有者"`                   // size:42: 地址长度限制 | comment: 数据库字段注释
	Blockchain      string    `gorm:"size:20;default:'sepolia';comment:区块链网络"` // size:20: 网络名称长度 | default:'sepolia': 默认值为'sepolia' | comment: 数据库字段注释
	LastSyncTime    time.Time `gorm:"comment:最后同步时间"`                          // comment: 数据库字段注释，记录最后同步时间
	CreatedAt       time.Time `gorm:"autoCreateTime"`                          // autoCreateTime: 自动设置创建时间
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`                          // autoUpdateTime: 自动更新时间
}

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名（如果未指定，GORM默认使用结构体名的复数形式）
func (Auction) TableName() string {
	return "auctions" // 明确指定表名为"auctions"，而不是GORM默认的复数形式"auctions"（这里相同，但习惯性显式声明）
}
