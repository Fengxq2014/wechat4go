package models

type GetConfigRespose struct {
	Openid    string `json:"openid"`
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
	Subscribe string `json:"subscribe"`
	Unionid   string `json:"unionid"`
}
