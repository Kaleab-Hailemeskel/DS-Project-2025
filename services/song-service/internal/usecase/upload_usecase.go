package usecase

import (
	"song-service/api/internal/domain"
)

type UploadUsecase struct {
	songRepo domain.ISongRepo
}

// UploadFileToArchive implements IUploadUsecase.
func (u *UploadUsecase) SaveSongMetaData(songMetaData *domain.Song) (*domain.Song, error) {
	//! checking for duplicate song in the database would be good before saving it
	return u.songRepo.SaveSong(songMetaData) //? save the song's metadata
}

func NewUploadUsecase(songRepo_ domain.ISongRepo) IUploadUsecase {
	return &UploadUsecase{
		songRepo: songRepo_,
	}
}
