package repositories

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileRepository interface {
	SaveProfile(fileHeader *multipart.FileHeader, userID string) (string, error)
}

type localFileRepository struct {
	basePath string
}

func NewLocalFileRepository(basePath string) FileRepository {

	path := "./files/profiles"
	os.MkdirAll(path, os.ModePerm)
	return &localFileRepository{basePath: path}
}

func (f *localFileRepository) SaveProfile(fileHeader *multipart.FileHeader, userID string) (string, error) {
	extension := filepath.Ext(fileHeader.Filename)
	fileName := userID + extension

	diskPath := filepath.Join(f.basePath, fileName)

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(diskPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	publicPath := "/files/profiles/" + fileName
	return publicPath, nil
}
