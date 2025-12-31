package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"song-service/api/config"
	"song-service/api/internal/domain"
	"song-service/api/internal/usecase"
	"song-service/api/pkg/media"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	songUsecase usecase.ISongUsecase
	media       media.IMediaProcessor
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

// GetAllSongs implements [IController].
func (u *Controller) GetAllSongs(ctx *gin.Context) {
	pageLimit := ctx.DefaultQuery("page-limit", fmt.Sprint(config.MAX_PAGE_SIZE))
	pageNumber := ctx.DefaultQuery("page-number", "1")
	songs, err := u.songUsecase.GetAllSong(pageNumber, pageLimit)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": "Failed to fetch songs: " + err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"songs": songs})
}

// UploadFileToArchive implements IController.
func (u *Controller) UploadFileToArchive(ctx *gin.Context) {
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
	songMetadata, err := u.songUsecase.SaveSongMetaData(&metadata)
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
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to save file: " + err.Error()},
		)
		return
	}

	// Use buffered channels (size 1) to prevent goroutines from hanging
	albumExtError := make(chan error, 1)
	segError := make(chan error, 1)
	blobChan := make(chan []byte, 1)

	// Goroutine 1: Album Art Extraction
	go func() {
		value, errEx := u.media.ExtractAlbumArt(saveDir)
		log.Println("✅ Internal: Extraction Finished")

		// IMPORTANT: Send the error first because the main thread reads it first
		albumExtError <- errEx
		blobChan <- value
	}()

	// Goroutine 2: HLS Segmentation
	go func() {
		log.Println("✅ Internal: Starting HLS Segments")
		errSeg := u.media.CreateHLSSegments(saveDir, parentFolder, "segment_%03d.ts", "index.m3u8")
		segError <- errSeg
	}()

	// --- RECEIVE SECTION ---

	// 1. Check HLS error
	segErr := <-segError
	if segErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to chunk file: " + segErr.Error()})
		return
	}

	// 2. Check Album Art error
	albumErr := <-albumExtError
	if albumErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get album art: " + albumErr.Error()})
		return
	}

	// 3. Finally, get the blob (now guaranteed to be available)
	albumArt := <-blobChan
	log.Println("✅ Success: Received blob of size", len(albumArt))

	// Save to Database
	err = u.songUsecase.SaveBlobImage(songMetadata.ID, albumArt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save album art: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "music created successfully"})
}

// SearchSongs implements IController.
func (s *Controller) SearchSongs(ctx *gin.Context) {
	titlePrefix := ctx.Query("title-prefix")
	titlePrefix = strings.TrimSpace(titlePrefix)
	pageLimit := ctx.DefaultQuery("page-limit", fmt.Sprint(config.MAX_PAGE_SIZE))
	pageNumber := ctx.DefaultQuery("page-number", "1")
	if titlePrefix == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Missing required query parameter: title-prefix"},
		)
		return
	}
	log.Println("✅ TitlePrefix =>", titlePrefix)
	// Call the usecase to search songs
	songs, err := s.songUsecase.SearchSongsByPrefix(titlePrefix, pageNumber, pageLimit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to search songs: " + err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"songs": songs})
}

func NewController(songUsecase_ usecase.ISongUsecase, media_ media.IMediaProcessor) IController {
	return &Controller{
		songUsecase: songUsecase_,
		media:       media_,
	}
}
