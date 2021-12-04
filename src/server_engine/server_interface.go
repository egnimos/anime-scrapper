package server_engine

import (
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/utility"
)

type AnimeServerInterface interface {
	AnimeListingSelector(serverCount int, pageCount int) map[string]interface{}
	AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings)
	AnimeInfoSelector(serverCount int, path string) map[string]interface{}
	AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo)
	EpisodesSelector(serverCount int, pageCount int, path string) map[string]interface{}
	EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{})
	EpisodesInfoSelector(serverCount int, path string) map[string]interface{}
	EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string)
}

// type ChoosenServerInterface interface {
// 	AnimeServerInterface
// }
