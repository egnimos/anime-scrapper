package nine_servers

import (
	"fmt"

	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetNineAnimeGG server_engine.AnimeServerInterface = &nineAnimeGG{}
)

type nineAnimeGG struct{}

func (n *nineAnimeGG) getServerKey(serverCount int) string {
	serverName := "9_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

//get anime listing selectors
func (n *nineAnimeGG) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeGG) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	return nil, nil
}

func (n *nineAnimeGG) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeGG) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	return nil, nil
}

func (n *nineAnimeGG) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeGG) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeGG) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeGG) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	return nil, nil
}

func (n *nineAnimeGG) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeGG) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}