package models

type GetConfigRespose struct {
	Openid    string `json:"openid"`
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
	Subscribe int32  `json:"subscribe"`
	Unionid   string `json:"unionid"`
}
