package usecase

import "song-service/api/internal/domain"


type IUploadUsecase interface {
	SaveSongMetaData(songMetaData *domain.Song) (*domain.Song, error)
}
