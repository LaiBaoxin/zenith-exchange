package request

type WithdrawRequest struct {
	User      string `json:"user" binding:"required"`
	Token     string `json:"token" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
	Nonce     int64  `json:"nonce"`
	VaultAddr string `json:"vault_addr" binding:"required"`
	ChainID   int64  `json:"chain_id" binding:"required"`
}
