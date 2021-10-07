package image

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nfnt/resize"
)

type Image struct {
	image   image.Image
	outType string
}

func NewImage(reader io.Reader, contentType string) (*Image, error) {
	var img image.Image
	var err error

	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(reader)
	case "image/png":
		img, err = png.Decode(reader)
	default:
		err = errors.New(fmt.Sprintf("the mime %s is not supported", contentType))
	}
	if err != nil {
		return nil, err
	}

	return &Image{
		image:   img,
		outType: "image/jpeg",
	}, nil
}

// Resize 修改为指定大小
func (r *Image) Resize(width, height uint, writer io.Writer) error {
	img := resize.Resize(width, height, r.image, resize.Lanczos3)
	return jpeg.Encode(writer, img, nil)
}

// Resizes 按预设将图片修改为指定大小
func (r *Image) Resizes(preset []*[2]uint, writers []bytes.Buffer) error {
	for i, writer := range writers {
		p := preset[i]
		err := r.Resize(p[0], p[1], &writer)
		if err != nil {
			return err
		}
	}

	return nil
}

// Thumbnail 放缩至预设最大宽高范围内
func (r *Image) Thumbnail(maxWidth, maxHeight uint, writer io.Writer) error {
	img := resize.Thumbnail(maxWidth, maxHeight, r.image, resize.Lanczos3)
	return jpeg.Encode(writer, img, nil)
}

// Thumbnails 批量放缩至预设最大宽高范围内
func (r *Image) Thumbnails(preset []*[2]uint, writers []io.Writer) error {
	for i, writer := range writers {
		p := preset[i]
		err := r.Thumbnail(p[0], p[1], writer)
		if err != nil {
			return err
		}
	}

	return nil
}

// Watermark 添加水印
func (r *Image) Watermark() error {
	return nil
}
