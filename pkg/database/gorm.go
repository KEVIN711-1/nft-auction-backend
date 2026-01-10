package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"nft-auction-backend/internal/config"
	"nft-auction-backend/internal/model"
)

func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// 1. 确保数据目录存在
	dataDir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %v", err)
	}

	// 2. 连接SQLite数据库
	log.Printf("连接数据库: %s", cfg.Path)
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 3. 自动迁移表（必须包含所有表！）
	if err := autoMigrateTables(db); err != nil {
		return nil, fmt.Errorf("创建表失败: %v", err)
	}

	log.Println("✅ 数据库初始化完成")
	return db, nil
}

func autoMigrateTables(db *gorm.DB) error {
	log.Println("开始自动创建表...")

	// 注册所有需要创建的表模型
	models := []interface{}{
		&model.Auction{}, // 拍卖表
		&model.NFTInfo{}, // NFT信息表（新增）
		// 可以添加更多表...
	}

	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			// 获取表名
			stmt := &gorm.Statement{DB: db}
			stmt.Parse(m)
			return fmt.Errorf("创建表 '%s' 失败: %v", stmt.Schema.Table, err)
		}

		stmt := &gorm.Statement{DB: db}
		stmt.Parse(m)
		log.Printf("✓ 表 '%s' 已就绪", stmt.Schema.Table)
	}

	return nil
}
