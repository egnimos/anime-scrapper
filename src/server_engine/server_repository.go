package server_engine

import (
	"fmt"

	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	ServerRepo ServerRepository = &serverRepository{}
)

type ServerRepository interface {
	AnimeListing(animeServer AnimeServerInterface, serverCount int, pageCount int) (utility.RestError, repository.AnimeListings)
	AnimeInfo(animeServer AnimeServerInterface, serverCount int, path string) (utility.RestError, *repository.AnimeInfo)
	AnimeEpisodes(animeServer AnimeServerInterface, serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{})
	AnimeEpisodeInfo(animeServer AnimeServerInterface, serverCount int, path string) (utility.RestError, []string)
}

type serverRepository struct{}

//define function of the given func type
// type task func() chromedp.Tasks

//get the list of anime
func (s *serverRepository) AnimeListing(animeServer AnimeServerInterface, serverCount int, pageCount int) (utility.RestError, repository.AnimeListings) {
	//get the anime listing selectors and url in form of MAP
	value := animeServer.AnimeListingSelector(serverCount, pageCount)
	fmt.Println(value)

	//parse the return HTML and convert it into the struct
	parseErr, animeListings := animeServer.AnimeListingHtmlParser(value)
	if parseErr != nil {
		return parseErr, nil
	}

	return nil, animeListings
}

//get the anime info
func (s *serverRepository) AnimeInfo(animeServer AnimeServerInterface, serverCount int, path string) (utility.RestError, *repository.AnimeInfo) {
	//get the anime info
	value := animeServer.AnimeInfoSelector(serverCount, path)

	//parse the return HTML and convert it into struct
	return animeServer.AnimeInfoHtmlParser(value)
}

func (s *serverRepository) AnimeEpisodes(animeServer AnimeServerInterface, serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	//get the anime episodes
	value := animeServer.EpisodesSelector(serverCount, pageCount, path)

	//parse the return HTML and convert it into struct
	return animeServer.EpisodesHtmlParser(value)
}

func (s *serverRepository) AnimeEpisodeInfo(animeServer AnimeServerInterface, serverCount int, path string) (utility.RestError, []string) {
	//get the anime info
	value := animeServer.EpisodesInfoSelector(serverCount, path)

	//parse the return HTML and convert it into struct
	return animeServer.EpisodesInfoHtmlParser(value)
}
