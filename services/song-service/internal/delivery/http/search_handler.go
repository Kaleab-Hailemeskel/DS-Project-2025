package http

import (
	"song-service/api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	searchUsecase usecase.ISearchEngineUsecase
}

// SearchSongs implements ISearchController.
func (s *SearchController) SearchSongs(ctx *gin.Context) {
	titlePrefix := ctx.Query("title_prefix")
	offset := ctx.DefaultQuery("offset", "0")
	page := ctx.DefaultQuery("page", "1")

	// Call the usecase to search songs
	songs, err := s.searchUsecase.SearchSongsByTitlePrefix(titlePrefix, offset, page)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to search songs: " + err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"songs": songs})
}

func NewSearchController(searchUsecase_ usecase.ISearchEngineUsecase) ISearchController {
	return &SearchController{
		searchUsecase: searchUsecase_,
	}
}
