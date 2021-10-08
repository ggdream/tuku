package controller

import (
	"bytes"
	"fmt"
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

// GetImage 获取图片
func GetImage(fsIns fs.FS) gin.HandlerFunc {
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
		var form model.GetImageReq
		if err := c.ShouldBind(&form); err != nil {
			errno.Abort(c, errno.TypeParamsParsingErr)
			return
		}
		object := c.Param("object")

		var needResize bool
		var preset [2]uint

		if form.Width*form.Height != 0 {
			p := [2]uint{form.Width, form.Height}
			var found bool
			for _, v := range presets {
				if p == *v {
					found = true
					break
				}
			}

			if !found {
				needResize = true
				preset = p
			}
		} else if form.Width == 0 && form.Height == 0 {
			needResize = true
			preset = *presets[0]
		} else {
			needResize = true
			preset = [2]uint{form.Width, form.Height}
		}

		if needResize {
			// 现读出原图并生成指定尺寸图片
			data, err := fsIns.ReadFile(object, true)
			if err != nil {
				errno.Abort(c, errno.TypeFileOpenFailed)
				return
			}

			// 获取文件类型
			buffer := make([]byte, 512)
			copy(buffer, data[:512])
			contentType := http.DetectContentType(buffer)

			// if contentType != "image/jpeg" && contentType != "image/png" {
			// 	errno.Abort(c, errno.TypeFileFormatNotSupported)
			// 	return
			// }

			imager, err := image.NewImage(bytes.NewReader(data), contentType)
			if err != nil {
				errno.Abort(c, errno.TypeFileOpenFailed)
				return
			}

			buf := bytes.NewBuffer([]byte{})
			err = imager.Resize(form.Width, form.Height, buf)
			if err != nil {
				errno.Abort(c, errno.TypeFileCannotFind)
			} else {
				errno.Stream(c, int64(buf.Len()), "image/jpeg", buf)
			}
		} else {
			// 从已有预设尺寸中获取
			uri := fmt.Sprintf("%s.%d_%d.jpg", object, preset[0], preset[1])
			reader, size, contentType, err := fsIns.GetReader(uri, false)
			if err != nil {
				errno.Abort(c, errno.TypeFileCannotFind)
				return
			}
			errno.Stream(c, size, contentType, reader)
		}
	}
}
