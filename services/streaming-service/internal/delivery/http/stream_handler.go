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
// @Summary Get HLS manifest
// @Description Returns the HLS manifest (index.m3u8) for a given song ID.
// @Tags streaming
// @Param filename path string true "Song ID"
// @Produce application/vnd.apple.mpegurl
// @Success 200 {string} string "m3u8 contents"
// @Failure 404 {object} map[string]string
// @Router /streams/get-manifest/{filename} [get]
func (s *StreamController) GetManifestFile(ctx *gin.Context) {
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

// GetStreamFilePath implements IStreamController.
// @Summary Get HLS segment
// @Description Returns a specific HLS segment (.ts) for a song.
// @Tags streaming
// @Param filename path string true "Song ID"
// @Param segment path string true "Segment filename"
// @Produce video/mp2t
// @Success 200 {string} string "binary segment"
// @Failure 400 {object} map[string]string
// @Router /streams/{filename}/{segment} [get]
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
