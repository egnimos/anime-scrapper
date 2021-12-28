package gogo_servers

/*
This source file contains the scraping of gogoanimeserver2.org
*/
import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/servers"

	// "github.com/egnimos/anime-scrapper/src/servers"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetGogoAnimeServer2Org server_engine.AnimeServerInterface = &gogoAnime2Org{}
)

type gogoAnime2Org struct{}

func (g *gogoAnime2Org) getServerKey(serverCount int) string {
	serverName := "gogo_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

// func (g *gogoAnime2Org) getAnimeListingUrl(serverKey string, pageCount int) string {
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

func (g *gogoAnime2Org) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerSu.AnimeListingSelector(serverCount, pageCount)
}

func (g *gogoAnime2Org) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	return GetGogoAnimeServerSu.AnimeListingHtmlParser(value)
}

//get the anime info
func (g *gogoAnime2Org) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerSu.AnimeInfoSelector(serverCount, path)
}

//anime info html parser
func (g *gogoAnime2Org) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	return GetGogoAnimeServerSu.AnimeInfoHtmlParser(value)
}

//get the lsit of episodes
func (g *gogoAnime2Org) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], path)
	// js := fmt.Sprintf(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_video_body ul#episode_page li a')[%d].click()`, pageCount)

	// var noOfPages int
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		//click to get the next set of episodes
		// chromedp.Evaluate(js, nil),
		//get the no of pages
		// chromedp.Evaluate(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_video_body ul#episode_page li a').length`, &noOfPages),
		chromedp.InnerHTML(`div#load_ep ul#episode_related`, &innerHTML),
	}

	//run the assigned task
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

//parse the episodes list
func (g *gogoAnime2Org) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	innerHTML := value["innerHTML"].(string)

	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	animeEpisodes := make([]repository.AnimeEpisode, 0)
	doc.Find(`li`).Each(func(i int, s *goquery.Selection) {
		episodesPathUrl := s.Find(`a`).AttrOr("href", "empty")
		episode := s.Find(`div.name`).Text()

		fmt.Printf("%s\n%s\n", episodesPathUrl, episode)
		animeEpisode := repository.AnimeEpisode{
			EpisodePath: episodesPathUrl,
			Episode:     episode,
		}

		animeEpisodes = append(animeEpisodes, animeEpisode)
	})

	return nil, map[string]interface{}{
		"pagination_status": "disabled",
		"episodes":          animeEpisodes,
	}
}

//get the episodes info
func (g *gogoAnime2Org) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
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
		"dataSource": fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], dataSrc),
	}
}

//parse the episodes info
func (g *gogoAnime2Org) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	dataSource := value["dataSource"].(string)
	dataSources := []string{dataSource}
	return nil, dataSources
}

func (g *gogoAnime2Org) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (g *gogoAnime2Org) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}
