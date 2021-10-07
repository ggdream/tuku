package errno

// getErrFromMap 从定义的errorsMap中获取文字描述信息
func getErrFromMap(t Type) string {
	return errorsMap[t]
}

var errorsMap = map[Type]string{
	TypePerfect: "",

	TypeFileOpenFailed:    "文件打开失败",
	TypeFileReadFailed:    "文件读取失败",
	TypeFileUploadFailed:  "文件上传失败",
	TypeFileCannotFind:    "文件找不到",
	TypeImageDecodeFailed: "图像解码失败",
	TypeImageResizeFailed: "图像调整失败",

	TypeParamsMissingErr: "参数缺失错误",
	TypeParamsParsingErr: "参数解析错误",
	TypeParamsInvalidErr: "参数无效错误",
}
