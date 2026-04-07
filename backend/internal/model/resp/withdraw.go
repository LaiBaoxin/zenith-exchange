package resp

type WithdrawResult struct {
	Signature string `json:"signature"`
	Nonce     uint64 `json:"nonce"`
	Amount    string `json:"amount"`
	Token     string `json:"token"`
}
