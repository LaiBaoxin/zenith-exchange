package utils

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

// DialEthClient 连接以太坊
func DialEthClient(url string) *ethclient.Client {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("无法连接到以太坊节点: %v", err)
	}
	return client
}
