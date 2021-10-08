package errno

const (
	// 正常
	TypePerfect = Type(iota)

	TypeFileOpenFailed         // 文件打开失败
	TypeFileReadFailed         // 文件读取失败
	TypeFileUploadFailed       // 文件上传失败
	TypeFileCannotFind         // 文件找不到
	TypeImageDecodeFailed      // 图像解码失败
	TypeImageResizeFailed      // 图像调整失败
	TypeFileFormatNotSupported // 不受支持的文件格式

	// 参数类
	TypeParamsMissingErr // 参数缺失错误
	TypeParamsParsingErr // 参数解析错误
	TypeParamsInvalidErr // 参数无效错误
)

type Type int16

func (t Type) Index() int16 {
	return int16(t)
}

func (t Type) String() string {
	return getErrFromMap(t)
}
