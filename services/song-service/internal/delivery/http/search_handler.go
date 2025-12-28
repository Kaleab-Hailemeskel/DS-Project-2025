package http

import (
	"fmt"
	"log"
	"net/http"
	"song-service/api/config"
	"song-service/api/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	searchUsecase usecase.ISearchEngineUsecase
}

// SearchSongs implements ISearchController.
func (s *SearchController) SearchSongs(ctx *gin.Context) {
	titlePrefix := ctx.Query("title-prefix")
	titlePrefix = strings.TrimSpace(titlePrefix)
	offset := ctx.DefaultQuery("page-limit", fmt.Sprint(config.MAX_PAGE_SIZE))
	page := ctx.DefaultQuery("page-number", "1")
	if titlePrefix == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Missing required query parameter: title-prefix"},
		)
		return
	}
	log.Println("✅ TitlePrefix =>", titlePrefix)
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
