package service

import (
	"fmt"
	"mime/multipart"

	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/google/cloudstorage"
	"github.com/pkg/errors"
)

// FileService ...
type FileService interface {
	UploadFile(bkt *cloudstorage.VBucket, roomID uint, fileUUID string, file *multipart.FileHeader, ftype FileType) (vItem *cloudstorage.VItem, err error)
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

func (ft *FileType) toString() string {
	var val string
	switch *ft {
	case AudioType:
		val = "audio"
	case PhotoType:
		val = "photo"
	}
	return val
}

// NewFileService ...
func NewFileService() FileService {
	return &fileService{}
}

// UploadFile uploads file to cloudstorage bucket
func (fs *fileService) UploadFile(bkt *cloudstorage.VBucket, roomID uint, fileUUID string, file *multipart.FileHeader, ftype FileType) (*cloudstorage.VItem, error) {
	openedFile, err := file.Open()
	if err != nil {
		return nil, err
	}

	path := fs.cloudStoragePath(roomID, fileUUID, ftype)
	var vItem *cloudstorage.VItem
	switch ftype {
	case AudioType:
		vItem, err = bkt.Put(path, openedFile, file.Size, make(map[string]interface{}))
	case PhotoType:
		vItem, err = bkt.Put(path, openedFile, file.Size, make(map[string]interface{}))
	default:
		return nil, errors.Errorf(`cannot upload file '%s' to storage server`, path)
	}

	if err != nil {
		return nil, err
	}
	return vItem, nil
}

func (fs *fileService) cloudStoragePath(roomID uint, fileUUID string, ftype FileType) string {
	return fmt.Sprintf(
		"%d/%s/%s",
		roomID,
		ftype.toString(),
		fileUUID,
	)
}
