package file

import (
	"errors"
	"os"
)

type Helper interface {
	CopyBytesArrayToDisk(dest string, datas []byte) error
	CreateDirectory(path string) error
	Read(path string) ([]byte, error)
	Delete(path string) error
}

type fileHelper struct{}

func NewHelper() *fileHelper {
	return &fileHelper{}
}

func (f *fileHelper) CopyBytesArrayToDisk(dest string, datas []byte) error {
	return os.WriteFile(dest, datas, 0666)
}

func (f *fileHelper) CreateDirectory(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return os.Mkdir("tmp", os.ModePerm)
	}

	return err
}

func (f *fileHelper) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (f *fileHelper) Delete(path string) error {
	return os.Remove(path)
}
