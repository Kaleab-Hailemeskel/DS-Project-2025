package usecase

import (
	"song-service/api/internal/domain"
	"song-service/api/internal/repository"
)

type UploadUsecase struct {
	songRepo repository.ISongRepo
}

// UploadFileToArchive implements IUploadUsecase.
func (u *UploadUsecase) SaveSongMetaData(songMetaData *domain.Song) (*domain.Song, error) {
	//! checking for duplicate song in the database would be good before saving it
	return u.songRepo.SaveSong(songMetaData) //? save the song's metadata
}

func NewUploadUsecase(songRepo_ repository.ISongRepo) IUploadUsecase {
	return &UploadUsecase{
		songRepo: songRepo_,
	}
}
