package gogo_servers

/*
This source file contains the scraping of gogoanimeserver.org
*/

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/servers"

	// "github.com/egnimos/anime-scrapper/src/servers"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetGogoAnimeServerOrg server_engine.AnimeServerInterface = &gogoAnimeOrg{}
)

type gogoAnimeOrg struct{}

func (g *gogoAnimeOrg) getServerKey(serverCount int) string {
	serverName := "gogo_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

// func (g *gogoAnimeOrg) getServerKey(serverCount int) string {
// 	serverName := "gogo_anime_server"
// 	return fmt.Sprintf("%s_%d", serverName, serverCount)
// }

// func (g *gogoAnimeOrg) getAnimeListingUrl(serverKey string, pageCount int) string {
// 	url := server_engine.ParsedServers.GogoanimeServers[serverKey]
// 	switch serverKey {
// 	case "gogo_anime_server_1": //https://gogoanime.cm
// 		return fmt.Sprintf("%s/anime-list.html?page=%d", url, pageCount)
// 	case "gogo_anime_server_2": //https://gogo-anime.su
// 		return fmt.Sprintf("%s/all_anime/?page=%d&alphabet=all", url, pageCount)
// 	case "gogo_anime_server_3": //https://ww2.gogoanimes.org
// 		return fmt.Sprintf("%s/anime-list?page=%d", url, pageCount)
// 	case "gogo_anime_server_4": //https://ww1.gogoanime2.org
// 		return fmt.Sprintf("%s/animelist/all/%d", url, pageCount)
// 	case "gogo_anime_server_5": //https://www1.gogoanime.sx
// 		return fmt.Sprintf("%s/anime-list?page=%d", url, pageCount)
// 	case "gogo_anime_server_6": //https://gogoanime.mom
// 		return fmt.Sprintf("%s/anime-list", url)
// 	case "gogo_anime_server_7": //https://gogoanime.wiki
// 		return fmt.Sprintf("%s/anime-list.html?page=%d", url, pageCount)
// 	default:
// 		return ""
// 	}
// }

func (g *gogoAnimeOrg) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.AnimeListingSelector(serverCount, pageCount)
}

func (g *gogoAnimeOrg) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	return GetGogoAnimeServerCm.AnimeListingHtmlParser(value)
}

//get the anime info
func (g *gogoAnimeOrg) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.AnimeInfoSelector(serverCount, path)
}

//anime info html parser
func (g *gogoAnimeOrg) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	return GetGogoAnimeServerCm.AnimeInfoHtmlParser(value)
}

//get the lsit of episodes
func (g *gogoAnimeOrg) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.EpisodesSelector(serverCount, pageCount, path)
}

//parse the episodes list
func (g *gogoAnimeOrg) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.EpisodesHtmlParser(value)
}

//get the episodes info
func (g *gogoAnimeOrg) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	// 
	serverKey := g.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], path)
	var dataSrc string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.AttributeValue(`div#load_anime div.anime_video_body_watch_items.load div.play-video iframe`, "src", &dataSrc, nil),
	}

	//run the tasks
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return nil, map[string]interface{}{
		"dataSource": dataSrc,
	}
}

//parse the episodes info
func (g *gogoAnimeOrg) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	dataSource := value["dataSource"].(string)
	dataSources := []string{dataSource}
	return nil, dataSources
}

func (g *gogoAnimeOrg) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	
	return nil, nil
}

func (g *gogoAnimeOrg) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}
