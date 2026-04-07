package main

import (
	"fmt"
	"log"

	"github.com/wwater/zenith-exchange/backend/internal/controller"
	"github.com/wwater/zenith-exchange/backend/internal/router" // 引入新路由包
	"github.com/wwater/zenith-exchange/backend/pkg/config"
)

func main() {
	// 加载配置（包括从 IO 读取私钥）
	config.InitConfig()

	// 初始化各层 Handler， 获取合约地址
	vaultHandler := controller.NewVaultHandler(config.GlobalConfig.Blockchain.VaultAddress)

	// 调用 Router 直接注入 Handler
	r := router.SetupRouter(vaultHandler)

	// 获取端口号
	port := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)

	// 启动服务
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
