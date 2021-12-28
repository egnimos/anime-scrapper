package nine_servers

import (
	"fmt"

	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetNineAnimeVc server_engine.AnimeServerInterface = &nineAnimeVc{}
)

type nineAnimeVc struct{}

func (n *nineAnimeVc) getServerKey(serverCount int) string {
	serverName := "9_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

//get anime listing selectors
func (n *nineAnimeVc) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeVc) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	return nil, nil
}

func (n *nineAnimeVc) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeVc) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	return nil, nil
}

func (n *nineAnimeVc) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeVc) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeVc) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeVc) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	return nil, nil
}

func (n *nineAnimeVc) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeVc) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}

