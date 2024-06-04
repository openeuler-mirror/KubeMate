package proto

type BaseResult struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"request_id"`
}
