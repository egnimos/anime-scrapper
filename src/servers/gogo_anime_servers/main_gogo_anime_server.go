package gogo_anime_servers

import (
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/servers/gogo_anime_servers/gogo_servers"
)

func Run(serverCount int) server_engine.AnimeServerInterface {
	switch serverCount {
	case 1:
		return gogo_servers.GetGogoAnimeServerCm
	case 2:
		return gogo_servers.GetGogoAnimeServerSu
	case 3:
		return gogo_servers.GetGogoAnimeServerOrg
	case 4:
		return gogo_servers.GetGogoAnimeServer2Org
	case 5:
		return gogo_servers.GetGogoAnimeServerSx
	case 6:
		return gogo_servers.GetGogoAnimeServerMom
	case 7:
		return gogo_servers.GetGogoAnimeServerWiki
	default:
		return nil
	}
}
