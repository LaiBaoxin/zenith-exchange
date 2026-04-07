package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/redis/go-redis/v9"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	CH  driver.Conn
	RDB *redis.Client // 全局 Redis 客户端
)

func InitDB() {
	var err error
	ctx := context.Background()

	// 初始化 MySQL
	dsn := config.GlobalConfig.DataBase.MySQL.Source
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("MySQL 连接失败: %v", err)
	}

	// 初始化 ClickHouse
	chCfg := config.GlobalConfig.DataBase.ClickHouse
	CH, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", chCfg.Host, chCfg.Port)},
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

	// 初始化 Redis
	rCfg := config.GlobalConfig.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rCfg.Host, rCfg.Port),
		Password: rCfg.Password,
		DB:       rCfg.DB,
		PoolSize: rCfg.PoolSize,
	})

	// 测试 Redis 连接
	if _, err = RDB.Ping(ctx).Result(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	// 自动迁移
	DB.AutoMigrate(&model.User{}, &model.Account{}, &model.Order{})

	log.Println("数据库环境就绪: MySQL, ClickHouse, Redis")
}
