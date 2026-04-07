package service

import (
	"github.com/wwater/zenith-exchange/backend/pkg/config"
)

type SystemService struct{}

type SystemConfigResp struct {
	VaultAddress string `json:"vault_address"`
	TokenAddress string `json:"token_address"`
	ChainID      uint64 `json:"chain_id"`
}

func (s *SystemService) GetGlobalConfig() SystemConfigResp {
	cfg := config.GlobalConfig.Blockchain
	return SystemConfigResp{
		VaultAddress: cfg.VaultAddress,
		TokenAddress: cfg.TokenAddress,
		ChainID:      cfg.ChainID,
	}
}
