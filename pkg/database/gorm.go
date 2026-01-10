// database 包提供数据库初始化和管理的功能
package database

import (
	"fmt"           // 格式化输入输出，用于创建格式化错误信息
	"log"           // 日志记录，用于输出程序运行状态
	"os"            // 操作系统功能，用于目录和文件操作
	"path/filepath" // 路径处理，用于处理文件路径分隔符和目录操作

	"github.com/glebarez/sqlite" // SQLite 数据库驱动（纯Go实现）
	"gorm.io/gorm"               // ORM框架主包
	"gorm.io/gorm/logger"        // GORM的日志器

	"nft-auction-backend/internal/config" // 项目配置模块
	"nft-auction-backend/internal/model"  // 数据模型定义
)

// InitDB 初始化数据库连接并创建表结构
// 参数: cfg - 数据库配置，包含数据库文件路径等信息
// 返回值: *gorm.DB - 数据库连接对象
//
//	error   - 错误信息，成功时为nil
func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// 1. 确保数据目录存在
	// filepath.Dir() 获取路径的目录部分，例如 "./data/db.sqlite" → "./data"
	dataDir := filepath.Dir(cfg.Path)

	// os.MkdirAll() 递归创建目录，如果目录已存在则什么都不做
	// 0755 是Unix权限位：所有者可读可写可执行，组和其他用户可读可执行
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		// fmt.Errorf() 创建格式化错误，包含原始错误信息
		return nil, fmt.Errorf("创建数据目录失败: %v", err)
	}

	// 2. 连接SQLite数据库
	// sqlite.Open() 使用SQLite驱动打开数据库文件
	// gorm.Open() 建立GORM数据库连接，初始化gorm 配置信息
	log.Printf("连接数据库: %s", cfg.Path)
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		// 配置GORM日志器，设置日志级别为Info（显示SQL语句等信息）
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 3. 自动迁移表（必须包含所有表！）
	// autoMigrateTables() 根据模型定义自动创建或更新表结构
	if err := autoMigrateTables(db); err != nil {
		return nil, fmt.Errorf("创建表失败: %v", err)
	}

	// 输出成功日志
	log.Println("✅ 数据库初始化完成")

	// 返回数据库连接对象
	return db, nil
}

// autoMigrateTables 自动创建或更新数据库表结构
// 参数: db - 已连接的数据库对象
// 返回值: error - 错误信息，成功时为nil
func autoMigrateTables(db *gorm.DB) error {
	log.Println("开始自动创建表...")

	// 注册所有需要创建的表模型
	// 使用interface{}类型切片，可以存放任意类型的模型指针
	models := []interface{}{
		&model.Auction{}, // 拍卖表模型
		&model.NFTInfo{}, // NFT信息表模型
		// 可以添加更多表模型...
	}

	// 遍历所有模型，逐个创建表
	for _, m := range models {
		// db.AutoMigrate() 自动迁移：
		// 1. 如果表不存在，创建新表
		// 2. 如果表存在但结构有变化，更新表结构（添加缺失的列）
		// 注意：不会删除列或修改列类型
		if err := db.AutoMigrate(m); err != nil {
			// 获取表名用于错误信息
			// stmt 是GORM语句对象，用于解析模型信息
			stmt := &gorm.Statement{DB: db}
			stmt.Parse(m) // 解析模型，获取表名等信息

			return fmt.Errorf("创建表 '%s' 失败: %v", stmt.Schema.Table, err)
		}

		// 再次解析模型，获取表名用于日志输出
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(m)
		log.Printf("✓ 表 '%s' 已就绪", stmt.Schema.Table)
	}

	// 所有表创建成功，返回nil
	return nil
}

// 补充说明：
// 1. GORM的AutoMigrate功能：
//    - 自动创建主键、索引
//    - 自动设置字段类型（根据Go类型推断）
//    - 自动添加created_at、updated_at、deleted_at时间戳（如果模型包含gorm.Model）
//
// 2. SQLite驱动特点：
//    - github.com/glebarez/sqlite 是纯Go实现，无需CGO
//    - 支持交叉编译
//    - 性能较好
//
// 3. 日志级别说明：
//    - logger.Silent: 静默模式，不输出任何日志
//    - logger.Error: 只输出错误日志
//    - logger.Warn: 输出警告和错误
//    - logger.Info: 输出信息、警告、错误（包含SQL语句）
//
// 4. 注意事项：
//    - 在生产环境中，建议将日志级别调高（如logger.Warn或logger.Error）
//    - AutoMigrate不会处理数据迁移，复杂结构变化需要手动处理
//    - 建议在开发环境使用AutoMigrate，生产环境使用迁移脚本
