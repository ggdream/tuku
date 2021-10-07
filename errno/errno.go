package errno

import (
	"fmt"
	"github.com/ggdream/tuku/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"

)

// New 创建返回
// message 其中的第一个参数必须为string类型
func New(ctx *gin.Context, code Type, data interface{}, message ...interface{}) {
	jsonObj := model.Model{
		Code: code.Index(),
		Data: data,
	}

	if message == nil {
		jsonObj.Message = code.String()
	} else if len(message) == 1 {
		jsonObj.Message = message[0].(string)
	} else {
		jsonObj.Message = fmt.Sprintf(message[0].(string), message[1:]...)
	}

	ctx.JSON(http.StatusOK, jsonObj)
}

// Perfect 成功响应
func Perfect(ctx *gin.Context, data interface{}, message ...interface{}) {
	New(ctx, TypePerfect, data, message...)
}

// Abort 中止下行，直接响应
func Abort(ctx *gin.Context, code Type) {
	jsonObj := model.Model{
		Code:    code.Index(),
		Message: code.String(),
	}
	ctx.AbortWithStatusJSON(http.StatusOK, jsonObj)
}

func Stream(ctx *gin.Context, contentLength int64, contentType string, reader io.Reader)  {
	ctx.DataFromReader(http.StatusOK, contentLength, contentType, reader, nil)
}
