package http

import (
	"fmt"
	"song-service/api/config"
	"song-service/api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	searchUsecase usecase.ISearchEngineUsecase
}

// SearchSongs implements ISearchController.
// @Summary Search songs by title prefix
// @Description Returns paginated songs filtered by title prefix.
// @Tags search
// @Produce json
// @Param title_prefix query string false "Title prefix filter"
// @Param page-limit query int false "Page size" default(10)
// @Param page-number query int false "Page number" default(1)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /search/songs [get]
func (s *SearchController) SearchSongs(ctx *gin.Context) {
	titlePrefix := ctx.Query("title_prefix")
	offset := ctx.DefaultQuery("page-limit", fmt.Sprint(config.MAX_PAGE_SIZE))
	page := ctx.DefaultQuery("page-number", "1")

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
