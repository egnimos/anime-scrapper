package kickass_anime_servers

import (
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/servers/kickass_anime_servers/kickass_servers"
)

func Run(serverCount int) server_engine.AnimeServerInterface {
	switch serverCount {
	case 1:
		return kickass_servers.GetKickAssServerRo
	default:
		return nil
	}
}