package model

// Model 返回模型的外部包装体
type Model struct {
	Code    int16         `json:"code" form:"code"`
	Data    interface{} `json:"data" form:"data"`
	Message string      `json:"message" form:"message"`
}
