package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/monitor"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"log"
	"os"

	"github.com/wwater/zenith-exchange/backend/internal/controller"
	DB "github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/router"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
)

func main() {
	pwd, _ := os.Getwd()
	fmt.Println("当前程序执行路径:", pwd)
	// 加载配置（包括从 IO 读取私钥）
	config.InitConfig()

	// 设置gin的运行模式
	if config.GlobalConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化DB连接
	DB.InitDB()

	// 初始化hub
	hub := service.NewHub()
	go hub.Run()

	// 监听用户存款websocket
	monitorSvc := monitor.NewDepositMonitor(
		"ws://127.0.0.1:8545",
		config.GlobalConfig.Blockchain.VaultAddress,
		hub,
	)
	log.Printf("Starting Vault Monitor on %s", config.GlobalConfig.Blockchain.VaultAddress)
	go monitorSvc.Start()

	// 初始化各层 Handler， 获取合约地址
	vaultHandler := controller.NewVaultHandler(config.GlobalConfig.Blockchain.VaultAddress)
	authHandler := &controller.AuthHandler{}
	sysHandler := &controller.SystemHandler{}
	assetHandler := &controller.AssetsHandler{}

	// 调用 Router 直接注入 Handler
	r := router.SetupRouter(vaultHandler, authHandler, sysHandler, assetHandler)

	// 获取端口号
	port := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)

	// 启动服务
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
