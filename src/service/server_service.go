package service

import (
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/servers"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetServices Services = &services{}
)

type Services interface {
	GetAnimeListing(query string, serverCount int, pageCount int) (utility.RestError, repository.AnimeListings)
	GetAnimeInfo(server string, serverCount int, path string) (utility.RestError, *repository.AnimeInfo)
	GetAnimeEpisodes(server string, serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{})
	GetAnimeEpisodeInfo(server string, serverCount int, path string) (utility.RestError, []string)
}

type services struct{}

func (s *services) getAnimeInterface(server string) server_engine.AnimeServerInterface {
	switch server {
	case "kickass":
		return servers.GetKickAssServer
	case "9anime":
		return servers.GetNineAnimeServers
	default:
		return nil
	}
}

//get the list of anime
func (s *services) GetAnimeListing(server string, serverCount int, pageCount int) (utility.RestError, repository.AnimeListings) {
	//get the server interface
	animeServerInterface := s.getAnimeInterface(server)
	//pass it to the server repository and get the map value
	return server_engine.ServerRepo.AnimeListing(animeServerInterface, serverCount, pageCount)
}

//get the anime info
func (s *services) GetAnimeInfo(server string, serverCount int, path string) (utility.RestError, *repository.AnimeInfo) {
	//get the server interface
	animeServerInterface := s.getAnimeInterface(server)
	//pass it to the server repository and return the value
	return server_engine.ServerRepo.AnimeInfo(animeServerInterface, serverCount, path)
}

//get the anime episodes
func (s *services) GetAnimeEpisodes(server string, serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	//get the server interface
	animeServerInterface := s.getAnimeInterface(server)
	//pass it to the server repository and return the value
	return server_engine.ServerRepo.AnimeEpisodes(animeServerInterface, serverCount, pageCount, path)
}

//get the anime episode info4
func (s *services) GetAnimeEpisodeInfo(server string, serverCount int, path string) (utility.RestError, []string) {
	//get the server interface
	animeServerInterface := s.getAnimeInterface(server)
	//pass it to the server repository and return the value
	return server_engine.ServerRepo.AnimeEpisodeInfo(animeServerInterface, serverCount, path)
}
