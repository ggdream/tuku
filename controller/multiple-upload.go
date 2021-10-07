package controller

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ggdream/tuku/config"
	"github.com/ggdream/tuku/errno"
	"github.com/ggdream/tuku/lib/fs"
	"github.com/ggdream/tuku/lib/image"
	"github.com/ggdream/tuku/model"
)

// MultipleUpload 同时上传多个文件
func MultipleUpload(fsIns fs.FS) gin.HandlerFunc {
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
		var form model.MultipleUploadReq
		if err := c.ShouldBind(&form); err != nil {
			errno.Abort(c, errno.TypeParamsParsingErr)
			return
		}

		files := form.Files
		results := make([]*fs.StorageResult, len(files))
		for i, file := range files {
			r, err := file.Open()
			if err != nil {
				errno.Abort(c, errno.TypeFileOpenFailed)
				return
			}

			buf := bytes.NewBuffer([]byte{})
			size, err := io.Copy(buf, r)
			if err != nil || size != file.Size {
				errno.Abort(c, errno.TypeFileReadFailed)
				return
			}

			filename := file.Filename
			if form.Names[i] != "" {
				filename = form.Names[i]
			}
			// 存储原图
			result, err := fsIns.WriteFile(r, filename, file.Size, true)
			if err != nil {
				errno.Abort(c, errno.TypeFileUploadFailed)
				return
			}
			results = append(results, result)

			imager, err := image.NewImage(buf, file.Header.Get("Content-Type"))
			if err != nil {
				errno.Abort(c, errno.TypeImageDecodeFailed)
				return
			}

			writers := make([]bytes.Buffer, len(presets))
			err = imager.Resizes(presets, writers)
			if err != nil {
				errno.Abort(c, errno.TypeImageResizeFailed)
				// TODO: 删除图片的操作
				return
			}

			for idx, reader := range writers {
				name := fmt.Sprintf("%s.%d_%d.jpg", filename, presets[idx][0], presets[idx][1])
				_, err = fsIns.WriteFile(&reader, name, int64(reader.Len()), false)
				if err != nil {
					errno.Abort(c, errno.TypeFileUploadFailed)
					return
				}
			}
		}

		errno.Perfect(c, results)
	}
}
