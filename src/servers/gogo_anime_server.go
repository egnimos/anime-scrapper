package servers

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetGogoAnimeServers server_engine.AnimeServerInterface = &gogoAnime{}
)

type gogoAnime struct{}

func (g *gogoAnime) getServerKey(serverCount int) string {
	serverName := "gogo_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

func (g *gogoAnime) AnimeListingSelector(serverCount int, pageCount int) map[string]interface{} {
	serverKey := g.getServerKey(serverCount)
	url := server_engine.ParsedServers.GogoanimeServers[serverKey]
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
	}
	//run the assigned task
	InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

func (g *gogoAnime) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	return nil, nil
}

//get the anime info
func (g *gogoAnime) AnimeInfoSelector(serverCount int, path string) map[string]interface{} {
	return nil
}

//anime info html parser
func (g *gogoAnime) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	return nil, nil
}

//get the lsit of episodes
func (g *gogoAnime) EpisodesSelector(serverCount int, pageCount int, path string) map[string]interface{} {
	return nil
}

//parse the episodes list
func (g *gogoAnime) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

//get the episodes info
func (g *gogoAnime) EpisodesInfoSelector(serverCount int, path string) map[string]interface{} {
	return nil
}

//parse the episodes info
func (g *gogoAnime) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	return nil, nil
}

func (g *gogoAnime) SearchAnimeSelector(keyword string) map[string]interface{} {
	return nil
} 

func (g *gogoAnime) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}