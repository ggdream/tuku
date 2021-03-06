package fs

import "io"

const (
	RAW = "raw"
	GEN = "gen"
)

type FS interface {
	WriteFile(reader io.Reader, objectName, contentType string, fileSize int64, isRaw bool) (*StorageResult, error)
	GetReader(objectName string, isRaw bool) (io.Reader, int64, string, error)
	ReadFile(objectName string, isRaw bool) ([]byte, error)

	adjustName(objectName string) string
	getRawPath(objectName string) string
	getGenPath(objectName string) string
}

type StorageResult struct {
}
