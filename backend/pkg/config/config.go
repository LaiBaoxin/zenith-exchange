package config

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	}
	JWT struct {
		Secret     string `mapstructure:"secret"`
		ExpireHour int    `mapstructure:"expire_hour"`
	} `mapstructure:"jwt"`
	Blockchain struct {
		VaultAddress string `mapstructure:"vault_address"`
		TokenAddress string `mapstructure:"token_address"`
		KeyPath      string `mapstructure:"key_path"` // 存储私钥文件的路径
		ChainID      uint64 `mapstructure:"chain_id"`
	} `mapstructure:"blockchain"`
	DataBase struct {
		MySQL struct {
			Source string `mapstructure:"source"`
		} `mapstructure:"mysql"`
		ClickHouse struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			Database string `mapstructure:"database"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
		} `mapstructure:"clickhouse"`
	} `mapstructure:"database"`
}

var (
	GlobalConfig     *Config
	SignerPrivateKey string
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 优先级2
	viper.AddConfigPath("./internal/config")
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 本地文件读取
	loadPrivateKey()
}

func loadPrivateKey() {
	// 从配置中拿到文件路径
	path := GlobalConfig.Blockchain.KeyPath
	if path == "" {
		log.Fatal("未配置私钥路径 (key_path)")
	}

	// IO 读取文件内容
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("无法读取私钥文件 [%s]: %v", path, err)
	}

	// 去掉可能的空格或换行符
	key := strings.TrimSpace(string(content))
	if key == "" {
		log.Fatal("私钥文件内容为空")
	}

	SignerPrivateKey = key
}
