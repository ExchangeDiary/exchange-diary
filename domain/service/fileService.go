package service

import (
	"mime/multipart"

	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/google/cloudstorage"
	"github.com/pkg/errors"
)

// FileService ...
type FileService interface {
	UploadFile(bkt *cloudstorage.VBucket, fileName string, file *multipart.FileHeader, ftype FileType) (vItem *cloudstorage.VItem, err error)
}

type fileService struct{}

// FileType handles kind of static file type which voda handles.
type FileType int

const (
	// AudioType ...
	AudioType FileType = iota
	// PhotoType ...
	PhotoType
)

// NewFileService ...
func NewFileService() FileService {
	return &fileService{}
}

// UploadFile uploads file to cloudstorage bucket
func (fs *fileService) UploadFile(bkt *cloudstorage.VBucket, fileName string, file *multipart.FileHeader, ftype FileType) (*cloudstorage.VItem, error) {
	openedFile, err := file.Open()
	if err != nil {
		return nil, err
	}

	var vItem *cloudstorage.VItem
	switch ftype {
	case AudioType:
		vItem, err = bkt.Put(fileName, openedFile, file.Size, make(map[string]interface{}))
	case PhotoType:
		vItem, err = bkt.Put(fileName, openedFile, file.Size, make(map[string]interface{}))
	default:
		return nil, errors.Errorf(`cannot upload file '%s' to storage server`, fileName)
	}

	if err != nil {
		return nil, err
	}
	return vItem, nil
}
