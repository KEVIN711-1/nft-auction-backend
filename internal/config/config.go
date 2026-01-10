// internal/config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`     // 服务器配置
	Database   DatabaseConfig   `mapstructure:"database"`   // 数据库配置
	Blockchain BlockchainConfig `mapstructure:"blockchain"` // 区块链配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `mapstructure:"port"` // 服务器监听端口
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string `mapstructure:"path"` // SQLite数据库文件路径
}

// BlockchainConfig 区块链配置
type BlockchainConfig struct {
	RPCURL                 string `mapstructure:"rpc_url"`                  // 以太坊节点RPC URL
	NFTContractAddress     string `mapstructure:"nft_contract_address"`     // NFT合约地址
	AuctionContractAddress string `mapstructure:"auction_contract_address"` // 拍卖合约地址
}

// LoadConfig 加载配置文件
func LoadConfig() *Config {
	// 设置配置文件名称和类型
	viper.SetConfigName("config") // 配置文件名（不含扩展名）
	viper.SetConfigType("yaml")   // 配置文件类型

	// 添加配置文件搜索路径（按优先级顺序）
	viper.AddConfigPath(".")        // 当前目录
	viper.AddConfigPath("./config") // config目录

	// 设置默认值（当配置文件缺失或字段为空时使用）
	viper.SetDefault("server.port", 8080)                       // 默认端口8080
	viper.SetDefault("database.path", "./data/auctions.db")     // 默认数据库路径
	viper.SetDefault("blockchain.rpc_url", "")                  // 默认空RPC URL（演示模式）
	viper.SetDefault("blockchain.nft_contract_address", "")     // 默认空NFT合约地址
	viper.SetDefault("blockchain.auction_contract_address", "") // 默认空拍卖合约地址

	var cfg Config

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果找不到配置文件，使用默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("⚠️  未找到配置文件 config.yaml，使用默认配置")
		} else {
			// 找到文件但读取失败
			log.Printf("⚠️  读取配置文件失败: %v", err)
		}
	} else {
		// 成功读取配置文件
		log.Printf("✅ 加载配置文件: %s", viper.ConfigFileUsed())
	}

	// 将配置绑定到结构体
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("⚠️  解析配置文件失败: %v", err)
		// 继续使用默认值
	}

	// 打印加载的配置信息（用于调试）
	log.Printf("=== 加载的配置 ===")
	log.Printf("服务器端口: %d", cfg.Server.Port)
	log.Printf("数据库路径: %s", cfg.Database.Path)
	log.Printf("RPC URL: %s", cfg.Blockchain.RPCURL)
	log.Printf("NFT合约地址: %s", cfg.Blockchain.NFTContractAddress)
	log.Printf("拍卖合约地址: %s", cfg.Blockchain.AuctionContractAddress)
	log.Printf("RPC URL是否为空: %v", cfg.Blockchain.RPCURL == "")

	return &cfg
}
