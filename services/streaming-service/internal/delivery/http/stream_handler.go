package http

import (
	"net/http"
	"streaming-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type StreamController struct {
	streamUsecase usecase.IStreamUsecase
}

// GetManifaste implements IStreamController.
func (s *StreamController) GetManifasteFile(ctx *gin.Context) {
	songId := ctx.Param("filename")
	filePath, err := s.streamUsecase.GetStreamFilePath(songId, "index.m3u8")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "HLS Manifest file not found. Did you run the segmenter first?"})
		return
	}
	// Set the correct Content-Type for the HLS manifest
	ctx.Header("Content-Type", "application/vnd.apple.mpegurl")
	ctx.File(filePath)
}

func (s *StreamController) GetManifasteFileTrial(ctx *gin.Context) {
	songId := ctx.Param("filename")
	filePath, err := s.streamUsecase.GetStreamFilePathTrial("index.m3u8")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "HLS Manifest file not found. Did you run the segmenter first?"})
		return
	}
	// Set the correct Content-Type for the HLS manifest
	ctx.Header("Content-Type", "application/vnd.apple.mpegurl")
	ctx.File(filePath)
}

// GetStreamFilePath implements IStreamController.
func (s *StreamController) GetStreamFile(ctx *gin.Context) {
	songId := ctx.Param("filename")
	segment := ctx.Param("segment")
	filePath, err := s.streamUsecase.GetStreamFilePath(songId, segment)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.File(filePath)
}

func NewStreamController(streamUseCase_ usecase.IStreamUsecase) IStreamController {
	return &StreamController{
		streamUsecase: streamUseCase_,
	}
}
