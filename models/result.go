package models

// Result api返回model
type Result struct {
	Res  int         `json:"result"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
