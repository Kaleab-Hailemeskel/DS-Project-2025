package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"streaming-service/config"
)

type StreamUsecase struct {
	musicArchiveDir string
}

// GetStreamFilePath implements IStreamUsecase.
func (s *StreamUsecase) GetStreamFilePath(songId string, segmentPos string) (string, error) {
	filePath_ := filepath.Join(config.SONG_ARCHIVE_DIR, songId, segmentPos)

	// Check if the file exists locally
	if _, err := os.Stat(filePath_); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s in %s", songId+segmentPos, config.SONG_ARCHIVE_DIR)
	}
	return filePath_, nil
}

func NewStreamUsecase(musicArchiveDir_ string) IStreamUsecase {
	return &StreamUsecase{
		musicArchiveDir: musicArchiveDir_,
	}
}
