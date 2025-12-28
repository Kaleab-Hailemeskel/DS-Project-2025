package media

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"song-service/api/config"
)

type HLSSegmenter struct {
}

// CreateHLSSegments implements domian.IHLSChunker.
func (h *HLSSegmenter) CreateHLSSegments(inputFile, outputDir, segNameTemplate, playlistName string) error {
	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	playlistPath := filepath.Join(outputDir, playlistName)
	segmentPathTemplate := filepath.Join(outputDir, segNameTemplate)

	// FFmpeg command arguments for HLS audio-only segmentation
	args := []string{
		"-i", inputFile,
		"-c:a", "aac",
		"-b:a", "128k",
		"-vn",
		"-f", "hls",
		"-hls_time", fmt.Sprintf("%d", config.SEGMENT_DURATION), // Segment duration in seconds
		"-hls_list_size", "0", // List all segments (for VOD)
		"-hls_segment_type", "mpegts", // Use MPEG-TS segments
		"-hls_segment_filename", segmentPathTemplate,
		playlistPath,
	}

	fmt.Printf("\n--- Executing FFmpeg ---\n")
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %w (Is FFmpeg installed and in PATH?)", err)
	}

	fmt.Printf("FFmpeg successfully created HLS segments in: %s\n", outputDir)
	// uncomment this if the original file needs to be deleted after chunking
	// os.Remove(inputFile) //? delete the original file after chunking
	return nil
}

func NewHLSSegmenter() IHLSChunker {
	return &HLSSegmenter{}
}
