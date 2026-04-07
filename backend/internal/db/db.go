package db // 建议改为 db，与目录名 internal/db 保持一致

import (
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局变量，首字母大写以便外部包访问
var (
	DB *gorm.DB
	CH driver.Conn
)

func InitDB() {
	var err error

	// 初始化 MySQL
	dsn := config.GlobalConfig.DataBase.MySQL.Source
	if dsn == "" {
		log.Fatal("MySQL 配置为空，请检查 config.yaml")
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("MySQL 连接失败: %v", err)
	}

	// 初始化 ClickHouse
	chCfg := config.GlobalConfig.DataBase.ClickHouse
	chAddr := fmt.Sprintf("%s:%d", chCfg.Host, chCfg.Port)

	CH, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{chAddr},
		Auth: clickhouse.Auth{
			Database: chCfg.Database,
			Username: chCfg.Username,
			Password: chCfg.Password,
		},
		DialTimeout: time.Second * 30,
	})
	if err != nil {
		log.Fatalf("ClickHouse 连接失败: %v", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(
		&model.User{},
		&model.Account{},
		&model.Order{},
	)
	if err != nil {
		log.Fatalf("MySQL 自动迁移失败: %v", err)
	}

	log.Println("数据库初始化成功：MySQL & ClickHouse")
}
