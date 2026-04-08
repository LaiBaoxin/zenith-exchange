package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/controller"
	DB "github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/monitor"
	"github.com/wwater/zenith-exchange/backend/internal/router"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
)

func main() {
	pwd, _ := os.Getwd()
	fmt.Println("当前程序执行路径:", pwd)

	// 加载配置
	config.InitConfig()

	// 设置 Gin 模式
	if config.GlobalConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化持久层连接 (MySQL, ClickHouse, Redis 的初始化)
	DB.InitDB()

	// 初始化通信层 Hub
	hub := service.NewHub()
	go hub.Run()

	// 初始化核心 Service 层
	matchSvc := service.NewMatchService(hub)
	if err := matchSvc.InitOrderBook(); err != nil {
		log.Fatalf("无法初始化订单簿: %v", err)
	}

	// 其他基础 Service
	orderSvc := service.NewOrderService(matchSvc)
	authSvc := &service.AuthService{}
	sysSvc := &service.SystemService{}
	klineSvc := service.NewKlineService()
	assetsSvc := service.NewAssetsService()
	marketSvc := service.NewMarketService(klineSvc, matchSvc)

	// 启动后台异步任务
	go marketSvc.StartPriceSimulationMock("BTC_USDT") // 模拟买卖交易
	go klineSvc.RunAggregator()

	// 监听用户存款合约
	monitorSvc := monitor.NewDepositMonitor(
		"ws://127.0.0.1:8545",
		config.GlobalConfig.Blockchain.VaultAddress,
		hub,
	)
	log.Printf("Starting Vault Monitor on %s", config.GlobalConfig.Blockchain.VaultAddress)
	go monitorSvc.Start()

	// 实例化 Controller 层
	vaultHandler := controller.NewVaultHandler(config.GlobalConfig.Blockchain.VaultAddress)
	authHandler := controller.NewAuthHandler(authSvc)
	sysHandler := controller.NewSystemHandler(sysSvc)
	assetHandler := controller.NewAssetsHandler(assetsSvc)
	wsHandler := controller.NewWSHandler(hub)
	marketHandler := controller.NewMarketHandler(marketSvc)
	orderHandler := controller.NewOrderHandler(orderSvc)

	// 装配路由
	r := router.SetupRouter(
		vaultHandler,
		authHandler,
		sysHandler,
		assetHandler,
		wsHandler,
		marketHandler,
		orderHandler,
	)

	port := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)

	// 启动服务
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
