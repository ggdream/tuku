package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ggdream/tuku/config"
	"github.com/ggdream/tuku/errno"
	"github.com/ggdream/tuku/lib/fs"
	"github.com/ggdream/tuku/lib/image"
	"github.com/ggdream/tuku/model"
)

// SimpleUpload 上传单个文件
func SimpleUpload(fsIns fs.FS) gin.HandlerFunc {
	rawPresets := config.Get().Image.Preset
	presets := make([]*[2]uint, len(rawPresets))
	for i, preset := range rawPresets {
		wh := strings.Split(preset, "*")
		if len(wh) != 2 {
			panic("please entry truly image preset")
		}

		width, err := strconv.Atoi(wh[0])
		if err != nil {
			panic(err)
		}
		height, err := strconv.Atoi(wh[1])
		if err != nil {
			panic(err)
		}
		presets[i] = &[2]uint{uint(width), uint(height)}
	}

	return func(c *gin.Context) {
		var form model.SimpleUploadReq
		if err := c.ShouldBind(&form); err != nil {
			errno.Abort(c, errno.TypeParamsParsingErr)
			return
		}

		file, err := form.File.Open()
		if err != nil {
			errno.Abort(c, errno.TypeFileOpenFailed)
			return
		}
		data, err := io.ReadAll(file)
		if err != nil {
			errno.Abort(c, errno.TypeFileReadFailed)
			return
		}

		// 获取文件类型
		buffer := make([]byte, 512)
		copy(buffer, data[:512])
		contentType := http.DetectContentType(buffer)

		if contentType != "image/jpeg" && contentType != "image/png" {
			errno.Abort(c, errno.TypeFileFormatNotSupported)
			return
		}

		// 先解码，判断是否为可解码再编辑的图片
		imager, err := image.NewImage(bytes.NewReader(data), contentType)
		if err != nil {
			errno.Abort(c, errno.TypeImageDecodeFailed)
			return
		}

		filename := form.File.Filename
		if form.Name != "" {
			filename = form.Name
		}
		// 存储原图
		result, err := fsIns.WriteFile(bytes.NewReader(data), filename, contentType, form.File.Size, true)
		if err != nil {
			fmt.Println(err)
			errno.Abort(c, errno.TypeFileUploadFailed)
			return
		}

		buffers, err := imager.Resizes(presets)
		if err != nil {
			errno.Abort(c, errno.TypeImageResizeFailed)
			// TODO: 删除图片的操作
			return
		}

		for i, reader := range buffers {
			name := fmt.Sprintf("%s.%d_%d.jpg", filename, presets[i][0], presets[i][1])
			_, err = fsIns.WriteFile(&reader, name, "image/jpeg", int64(reader.Len()), false)
			if err != nil {
				errno.Abort(c, errno.TypeFileUploadFailed)
				return
			}
		}

		errno.Perfect(c, result)
	}
}
