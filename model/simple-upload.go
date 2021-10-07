package model

import "mime/multipart"

type SimpleUploadReq struct {
	Name	string					`json:"name" form:"name"`
	File	*multipart.FileHeader	`json:"file" form:"file" binding:"required"`
}

type SimpleUploadRes struct {

}
