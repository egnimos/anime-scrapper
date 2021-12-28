package nine_anime_servers

import (
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/servers/9_anime_servers/nine_servers"
)

func Run(serverCount int) server_engine.AnimeServerInterface {
	switch serverCount {
	case 1:
		return nine_servers.GetNineAnime2Com
	default:
		return nil
	}
}