package usecase

type IStreamUsecase interface {
	GetStreamFilePath(songId, segmentPos string) (string, error) // returns the file path and an error if the song with songID segmentPos doesn't exist
}
