package service

import (
	"errors"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wwater/zenith-exchange/backend/internal/contract"
	"github.com/wwater/zenith-exchange/backend/internal/db"
	"github.com/wwater/zenith-exchange/backend/internal/model"
	"github.com/wwater/zenith-exchange/backend/pkg/config"
	"github.com/wwater/zenith-exchange/backend/pkg/utils"
)

type VaultService struct {
	cfg   *config.Config
	vault *contract.ZenithVault
}

func NewVaultService(cfg *config.Config, rpcURL string) *VaultService {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("VaultService: 无法连接以太坊节点: %v", err)
	}

	vaultAddr := common.HexToAddress(cfg.Blockchain.VaultAddress)
	vaultContract, err := contract.NewZenithVault(vaultAddr, client)
	if err != nil {
		log.Fatalf("VaultService: 无法绑定金库合约: %v", err)
	}

	log.Printf("VaultService 初始化完成: Vault=%s", cfg.Blockchain.VaultAddress)

	return &VaultService{
		cfg:   cfg,
		vault: vaultContract,
	}
}

// PrepareWithdraw 校验余额、生成签名，不修改数据库。
// 链上交易确认后由 VaultMonitor 监听 Withdraw 事件进行最终扣款。
func (s *VaultService) PrepareWithdraw(userID int64, walletAddr, currency, amountStr string) (string, int64, error) {
	var account model.Account
	if err := db.DB.Where("user_id = ? AND currency = ?", userID, currency).First(&account).Error; err != nil {
		return "", 0, errors.New("账户不存在")
	}

	availableWei, _ := DecimalToWei(account.Available)
	withdrawWei, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return "", 0, errors.New("提现金额格式非法")
	}
	if availableWei.Cmp(withdrawWei) < 0 {
		return "", 0, errors.New("可用余额不足")
	}

	userAddress := common.HexToAddress(walletAddr)
	onChainNonce, err := s.vault.Nonces(&bind.CallOpts{}, userAddress)
	if err != nil {
		return "", 0, errors.New("无法从链上读取 nonce: " + err.Error())
	}
	nonce := onChainNonce.Int64()

	sig, err := utils.SignWithdrawData(
		walletAddr,
		s.cfg.Blockchain.TokenAddress,
		amountStr,
		nonce,
		s.cfg.Blockchain.VaultAddress,
		int64(s.cfg.Blockchain.ChainID),
	)
	if err != nil {
		return "", 0, errors.New("后端签名生成异常: " + err.Error())
	}

	return sig, nonce, nil
}

func DecimalToWei(decimalStr string) (*big.Int, bool) {
	parts := strings.Split(decimalStr, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}

	if len(decPart) > 18 {
		decPart = decPart[:18]
	} else {
		decPart = decPart + strings.Repeat("0", 18-len(decPart))
	}

	return new(big.Int).SetString(intPart+decPart, 10)
}

func WeiToDecimalStr(wei *big.Int) string {
	if wei == nil {
		return "0"
	}
	s := wei.String()
	if len(s) <= 18 {
		return "0." + strings.Repeat("0", 18-len(s)) + s
	}
	loc := len(s) - 18
	return s[:loc] + "." + s[loc:]
}
