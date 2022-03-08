package controller

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/ExchangeDiary/exchange-diary/application"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/google/cloudstorage"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const bucketName = "voda_bucket"

// FileController handles /v1/files api
type FileController interface {
	// Get() gin.HandlerFunc
	Post() gin.HandlerFunc
}

type fileController struct {
	fileService service.FileService
}

// NewFileController ...
func NewFileController(fs service.FileService) FileController {
	return &fileController{fileService: fs}
}

// // @Summary      upload a file
// // @Description  파일 업로드
// // @Tags         filess
// func (fc *fileController) Get() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.JSON(http.StatusOK, res)
// 	}
// }

type photoFileRequest struct {
	Photo *multipart.FileHeader `form:"photo" swaggerignore:"true"`
}
type audioFileRequest struct {
	Audio *multipart.FileHeader `form:"audio" swaggerignore:"true"`
}

type filePostRequestForm struct {
	PhotoUUID string `form:"photoUUID" description:"유니크한 uuid" example:"e4947e0c-490b-4588-a14d-e74dd3b8371f"`
	photoFileRequest
	AudioUUID string `form:"audioUUID" description:"유니크한 uuid" example:"ad5bb198-942f-4ddf-a248-3aaa4bba3b9b"`
	audioFileRequest
	AudioTitle string `form:"audioTitle" description:"오디오 파일명" example:"LastDayOnEarth"`
	AudioPitch string `form:"audioPitch" description:"오디오 피치값" example:"1.5"`
}

type filePostResponse struct {
	PhotoURL         string `json:"photoURL,omitempty"`
	AudioURL         string `form:"audioURL,omitempty"`
	PhotoDownloadURL string `form:"photoDownloadURL,omitempty"`
	AudioDownloadURL string `form:"audioDownloadURL,omitempty"`
}

// @Summary      upload static files
// @Description  파일 업로드
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Param        room_id    path      int   true  "Room ID"
// @Param        form  formData  filePostRequestForm  true  "form"
// @Param        photo  formData  file  false  "photo"
// @Param        audio  formData  file  false  "audio"
// @Success      200  {object}   filePostResponse
// @Failure      400
// @Failure      500
// @Router       /rooms/{room_id}/files [post]
// @Security ApiKeyAuth
func (fc *fileController) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var fileForm filePostRequestForm
		roomID, err := application.ParseUint(c.Param("room_id"))
		if err := c.MustBindWith(&fileForm, binding.FormMultipart); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}

		vc, bkt, err := fc.getStorageSession()
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}
		defer vc.Close()
		res := filePostResponse{}
		// audio check
		if fileForm.AudioUUID != "" {
			vItem, err := fc.fileService.UploadFile(bkt, roomID, fileForm.AudioUUID, fileForm.Audio, service.AudioType)
			if err != nil {
				logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
					"error":   true,
				})
				return
			}
			res.AudioURL = vItem.URL().String()
			res.AudioDownloadURL = vItem.DownloadURL().String()
		}
		// photo check
		if fileForm.PhotoUUID != "" {
			vItem, err := fc.fileService.UploadFile(bkt, roomID, fileForm.PhotoUUID, fileForm.Photo, service.PhotoType)
			if err != nil {
				logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
					"error":   true,
				})
				return
			}
			res.PhotoURL = vItem.URL().String()
			res.PhotoDownloadURL = vItem.DownloadURL().String()
		}
		c.JSON(http.StatusOK, res)
	}
}

func (fc *fileController) getStorageSession() (*cloudstorage.VClient, *cloudstorage.VBucket, error) {
	vc, err := cloudstorage.GetVClient(context.Background())
	if err != nil {
		return nil, nil, err
	}
	bkt, err := vc.VBucket(bucketName)
	if err != nil {
		return nil, nil, err
	}
	return vc, bkt, nil
}
