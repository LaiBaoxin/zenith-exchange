package main

import (
	"fmt"
	"log"

	"github.com/wwater/zenith-exchange/backend/internal/controller"
	DB "github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/router"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
)

func main() {
	// 加载配置（包括从 IO 读取私钥）
	config.InitConfig()

	// 初始化DB连接
	DB.InitDB()

	// 初始化各层 Handler， 获取合约地址
	vaultHandler := controller.NewVaultHandler(config.GlobalConfig.Blockchain.VaultAddress)
	authHandler := &controller.AuthHandler{}
	sysHandler := &controller.SystemHandler{}

	// 调用 Router 直接注入 Handler
	r := router.SetupRouter(vaultHandler, authHandler, sysHandler)

	// 获取端口号
	port := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)

	// 启动服务
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
