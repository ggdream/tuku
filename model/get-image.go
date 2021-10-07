package model

type GetImageReq struct {
	Width  uint `json:"width" form:"width"`
	Height uint `json:"height" form:"height"`
}
