package image

import (
	"os"
	"testing"
)

func TestNewImage_Resize(t *testing.T) {
	file, err := os.Open("test.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	imager, err := NewImage(file, "image/jpeg")
	if err != nil {
		panic(err)
	}

	out, err := os.Create("out.jpg")
	if err != nil {
		panic(err)
	}
	err = imager.Resize(100, 0, out)
	if err != nil {
		panic(err)
	}
}

func TestNewImage_Thumbnail(t *testing.T) {
	file, err := os.Open("test.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	imager, err := NewImage(file, "image/jpeg")
	if err != nil {
		panic(err)
	}

	out, err := os.Create("res.jpg")
	if err != nil {
		panic(err)
	}
	err = imager.Thumbnail(100, 100, out)
	if err != nil {
		panic(err)
	}
}
