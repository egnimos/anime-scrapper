package application

import (
	"github.com/egnimos/anime-scrapper/src/controllers"
	ping_ "github.com/egnimos/anime-scrapper/src/controllers/ping"
	"github.com/gin-gonic/gin"
)

func routers(r *gin.Engine) {
	test := r.Group("/server")
	{
		test.GET("/ping", ping_.Ping.Ping)
	}
	//version 1
	v1 := r.Group("/v1")
	{
		v1.GET("/anime-listing/:server/server_count/:server_count/page_count/:page_count", controllers.GetControllers.GetAnimeListing)
		v1.GET("/anime-info/:server/server_count/:server_count", controllers.GetControllers.GetAnimeInfo)
		v1.GET("/anime-episodes/:server/server_count/:server_count/page_count/:page_count", controllers.GetControllers.GetAnimeEpisodes)
		v1.GET("/anime-episode-info/:server/server_count/:server_count", controllers.GetControllers.GetAnimeEpisodeInfo)
		v1.GET("/search/:server/server_count/:server_count/keyword/:keyword")
		v1.GET("")
	}
}
