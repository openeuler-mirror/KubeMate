package proto

type NodeInfoResult struct {
	BaseResult
	Data interface{} `json:"data"`
}

type NodeInfoParam struct {
	Label map[string]string `json:"label"`
}
