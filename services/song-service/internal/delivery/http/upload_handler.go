package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"song-service/api/config"
	"song-service/api/internal/domain"
	"song-service/api/internal/usecase"
	"song-service/api/pkg/media"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	uploadUsecase usecase.IUploadUsecase
	chunker       media.IHLSChunker
}

// Helper to clean up filenames (e.g., remove spaces/special chars)
func sanitizeFilename(s string) string {
	// Simple example: replace spaces with underscores.
	// A real implementation would be more robust.
	var result []rune
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result = append(result, r)
		} else if r == ' ' {
			result = append(result, '_')
		}
	}
	return string(result)
}

// UploadFileToArchive implements IUploadController.
func (u *UploadController) UploadFileToArchive(ctx *gin.Context) {
	// 1. **Get the uploaded file**
	file, err := ctx.FormFile("musicFile") // Key from the client must be "musicFile"
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get music file from form: " + err.Error()})
		return
	}

	// 2. **Retrieve the JSON metadata string from the form**
	// The client sends the entire metadata object as a string under the key "metadata".
	metadataJSON := ctx.PostForm("metadata")
	if metadataJSON == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing required 'metadata' JSON field in form"})
		return
	}

	// 3. **Unmarshal the JSON string into the MusicMetadata struct**
	var metadata domain.Song
	if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata format (must be valid JSON): " + err.Error()})
		return
	}
	songMetadata, err := u.uploadUsecase.SaveSongMetaData(&metadata)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
		return
	}
	fileExtension := filepath.Ext(file.Filename)

	uniqueFilename := fmt.Sprintf("%s%s", sanitizeFilename(songMetadata.ID.String()), fileExtension)
	parentFolder := filepath.Join(config.SONG_ARCHIVE_DIR, songMetadata.ID.String())
	saveDir := filepath.Join(parentFolder, uniqueFilename)
	// 5. **Save the file to the server's local folder**
	// This uses Gin's convenience function to save the file handle.
	if err := ctx.SaveUploadedFile(file, saveDir); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
		return
	}
	//? chunk the music and save it to the existing folder and uncomment a code found in the implementation to delete the original song file
	err = u.chunker.CreateHLSSegments(saveDir, parentFolder, "segment_%03d.ts", "index.m3u8")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to chunk file: " + err.Error()})
		return
	}
}

func NewUploadController(uploadUsecase_ usecase.IUploadUsecase) IUploadController {
	return &UploadController{
		uploadUsecase: uploadUsecase_,
	}
}
