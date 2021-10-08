package fs

import (
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	path string
}

func NewLocal(path string) (*Local, error) {
	if err := mkDir(filepath.Join(path, RAW)); err != nil {
		return nil, err
	}
	if err := mkDir(filepath.Join(path, GEN)); err != nil {
		return nil, err
	}

	return &Local{
		path: path,
	}, nil
}

func mkDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, os.ModePerm)
		} else {
			return err
		}
	}

	return nil
}

func (l *Local) WriteFile(reader io.Reader, objectName, contentType string, fileSize int64, isRaw bool) (*StorageResult, error) {
	path := l.getPath(objectName, isRaw)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (l *Local) GetReader(objectName string, isRaw bool) (io.Reader, int64, string, error) {
	path := l.getPath(objectName, isRaw)
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, "", err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, 0, "", err
	}

	return file, info.Size(), "image/jpeg", nil
}

func (l *Local) ReadFile(objectName string, isRaw bool) ([]byte, error) {
	path := l.getPath(objectName, isRaw)
	return os.ReadFile(path)
}

func (l *Local) getPath(objectName string, isRaw bool) string {
	adName := l.adjustName(objectName)
	if isRaw {
		return l.getRawPath(adName)
	} else {
		return l.getGenPath(adName)
	}
}

func (l *Local) adjustName(objectName string) string {
	return objectName
}

func (l *Local) getRawPath(objectName string) string {
	return filepath.Join(l.path, RAW, objectName)
}
func (l *Local) getGenPath(objectName string) string {
	return filepath.Join(l.path, GEN, objectName)
}
