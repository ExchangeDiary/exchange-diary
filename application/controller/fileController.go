package controller

import (
	"mime/multipart"
	"net/http"

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

// FilePostRequestForm ...
type FilePostRequestForm struct {
	PhotoUUID string                `form:"photoUUID"`
	Photo     *multipart.FileHeader `form:"photo"`
	AudioUUID string                `form:"audioUUID"`
	Audio     *multipart.FileHeader `form:"audio"`
}

// @Summary      upload a file
// @Description  파일 업로드
// @Tags         filess
func (fc *fileController) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var fileForm FilePostRequestForm
		if err := c.MustBindWith(&fileForm, binding.FormMultipart); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}

		bkt, err := fc.getBucket()
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}
		// audio check
		if fileForm.AudioUUID != "" {
			vItem, err := fc.fileService.UploadFile(bkt, fileForm.AudioUUID, fileForm.Audio, service.AudioType)
			if err != nil {
				logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
					"error":   true,
				})
				return
			}
			logger.Info("Audio: " + vItem.URL().String())
		}
		// photo check
		if fileForm.PhotoUUID != "" {
			vItem, err := fc.fileService.UploadFile(bkt, fileForm.PhotoUUID, fileForm.Photo, service.PhotoType)
			if err != nil {
				logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
					"error":   true,
				})
				return
			}
			logger.Info("Photo: " + vItem.URL().String())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "file uploaded successfully",
		})
	}
}

func (fc *fileController) getBucket() (*cloudstorage.VBucket, error) {
	// refs: https://www.vompressor.com/gin4/
	bkt, err := cloudstorage.GetClient().VBucket(bucketName)
	if err != nil {
		return nil, err
	}
	return bkt, nil
}
