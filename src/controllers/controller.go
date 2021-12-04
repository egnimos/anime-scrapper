package controllers

import (
	"net/http"
	"strconv"

	"github.com/egnimos/anime-scrapper/src/service"
	"github.com/egnimos/anime-scrapper/src/utility"
	"github.com/gin-gonic/gin"
)

var (
	GetControllers Controllers = &controllers{}
)

type Controllers interface {
	GetAnimeListing(ctx *gin.Context)
	GetSearchResult(ctx *gin.Context)
	GetAnimeInfo(ctx *gin.Context)
	GetAnimeEpisodes(ctx *gin.Context)
	GetAnimeEpisodeInfo(ctx *gin.Context)
}

type controllers struct {
	server      string
	serverCount int64
	pageCount   int64
	path        string
}

func (c *controllers) getServerParameters(ctx *gin.Context, isPageCount bool) utility.RestError {
	var err error
	c.server = ctx.Param("server")
	c.serverCount, err = strconv.ParseInt(ctx.Param("server_count"), 10, 64)
	if err != nil {
		return utility.NewBadRequestError("server count should be a number")
	}
	if isPageCount {
		c.pageCount, err = strconv.ParseInt(ctx.Param("page_count"), 10, 64)
		if err != nil {
			return utility.NewBadRequestError("page count should be a number")
		}
	}
	c.path = ctx.Query("path")
	return nil
}

//get the list of anime
func (c *controllers) GetAnimeListing(ctx *gin.Context) {
	paramController := &controllers{}
	if paramErr := paramController.getServerParameters(ctx, true); paramErr != nil {
		ctx.JSON(int(paramErr.Status()), paramErr)
		return
	}

	err, result := service.GetServices.GetAnimeListing(paramController.server, int(paramController.serverCount), int(paramController.pageCount))
	if err != nil {
		ctx.JSON(int(err.Status()), err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

//get the anime info
func (c *controllers) GetAnimeInfo(ctx *gin.Context) {
	paramController := &controllers{}
	if paramErr := paramController.getServerParameters(ctx, false); paramErr != nil {
		ctx.JSON(int(paramErr.Status()), paramErr)
		return
	}

	err, result := service.GetServices.GetAnimeInfo(paramController.server, int(paramController.serverCount), paramController.path)
	if err != nil {
		ctx.JSON(int(err.Status()), err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

//get the list of episodes
func (c *controllers) GetAnimeEpisodes(ctx *gin.Context) {
	paramController := &controllers{}
	if paramErr := paramController.getServerParameters(ctx, true); paramErr != nil {
		ctx.JSON(int(paramErr.Status()), paramErr)
		return
	}

	err, result := service.GetServices.GetAnimeEpisodes(paramController.server, int(paramController.serverCount), int(paramController.pageCount), paramController.path)
	if err != nil {
		ctx.JSON(int(err.Status()), err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

//get the anime episode info
func (c *controllers) GetAnimeEpisodeInfo(ctx *gin.Context) {
	paramController := &controllers{}
	if paramErr := paramController.getServerParameters(ctx, false); paramErr != nil {
		ctx.JSON(int(paramErr.Status()), paramErr)
		return
	}

	err, result := service.GetServices.GetAnimeEpisodeInfo(paramController.server, int(paramController.serverCount), paramController.path)
	if err != nil {
		ctx.JSON(int(err.Status()), err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *controllers) GetSearchResult(ctx *gin.Context) {

}
