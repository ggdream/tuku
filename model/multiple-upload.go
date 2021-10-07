package model

import "mime/multipart"

type MultipleUploadReq struct {
	Names []string                `json:"names" form:"names"`
	Files []*multipart.FileHeader `json:"files" form:"files" binding:"required"`
}
