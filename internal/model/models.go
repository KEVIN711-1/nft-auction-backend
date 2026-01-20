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
	Status        string `gorm:"size:16"`
	AuctionID     uint64 `gorm:"uniqueIndex"` // uniqueIndex: 创建唯一索引，防止重复值
	NFTContract   string `gorm:"size:42"`     // size:42: 字符串最大长度42个字符（以太坊地址长度）
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
	ID uint64 `gorm:"primaryKey;autoIncrement"`
	// ContractAddress string    `gorm:"uniqueIndex;size:42;comment:NFT合约地址"`
	ContractAddress string    `gorm:"size:42;not null;index:idx_contract_token;comment:NFT合约地址" json:"contract_address"`
	TokenID         string    `gorm:"size:100;comment:Token ID"` // 新增TokenID字段
	Name            string    `gorm:"size:255;comment:NFT名称"`
	Symbol          string    `gorm:"size:50;comment:NFT符号"`
	Uri             string    `gorm:"size:50;comment:URI"`
	TotalSupply     string    `gorm:"type:varchar(100);comment:总供应量"`
	Owner           string    `gorm:"size:42;comment:合约所有者"`
	Blockchain      string    `gorm:"size:20;default:'sepolia';comment:区块链网络"`
	LastSyncTime    time.Time `gorm:"comment:最后同步时间"`
	ContractName    string    `gorm:"size:255;comment:合约名称"`       // 新增
	ContractSymbol  string    `gorm:"size:50;comment:合约符号"`        // 新增
	IsMinted        bool      `gorm:"default:false;comment:是否已铸造"` // 新增
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
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

// EventSync 事件同步记录
type EventSync struct {
	ID        uint   `gorm:"primarykey"`
	EventType string `gorm:"size:50;uniqueIndex"` // auction_events, bid_events, nft_events
	LastBlock uint64 // 最后处理的区块
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Bid 当前出价记录（用于快速查询）
type Bid struct {
	ID          uint      `gorm:"primarykey"`
	AuctionID   uint64    `gorm:"index;not null"`             // 拍卖ID，加索引
	Bidder      string    `gorm:"size:42;not null;index"`     // 出价者地址
	Amount      string    `gorm:"type:varchar(100);not null"` // 出价金额
	TxHash      string    `gorm:"size:66;uniqueIndex"`        // 交易哈希
	Status      string    `gorm:"size:20;default:'pending'"`  // 状态: pending, confirmed, failed
	IsHighest   bool      `gorm:"default:false"`              // 是否是当前最高出价
	BlockNumber uint64    `gorm:"index"`                      // 区块高度
	ConfirmedAt time.Time // 确认时间
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// 更新 BidHistory，添加更多字段
type BidHistory struct {
	ID            uint      `gorm:"primarykey"`
	AuctionID     uint64    `gorm:"index"`                       // 拍卖ID
	Bidder        string    `gorm:"size:42"`                     // 出价者地址
	Amount        string    `gorm:"type:varchar(100)"`           // 出价金额
	TxHash        string    `gorm:"size:66;uniqueIndex"`         // 交易哈希
	Status        string    `gorm:"size:20;default:'submitted'"` // 状态: submitted, pending, confirmed, failed
	BlockNumber   uint64    `gorm:"index"`                       // 区块高度
	BlockTime     uint64    // 区块时间戳
	GasPrice      string    `gorm:"type:varchar(50)"` // Gas价格
	GasUsed       uint64    // Gas使用量
	Confirmations uint      `gorm:"default:0"` // 确认数
	ErrorMessage  string    `gorm:"type:text"` // 错误信息（如果有）
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
