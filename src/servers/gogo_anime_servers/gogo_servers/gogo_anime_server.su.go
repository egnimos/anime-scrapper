package gogo_servers

/*
This source file contains the scraping of gogoanimeserver.su
*/

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/servers"

	// "github.com/egnimos/anime-scrapper/src/servers"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetGogoAnimeServerSu server_engine.AnimeServerInterface = &gogoAnimeSu{}
)

type gogoAnimeSu struct{}

func (g *gogoAnimeSu) getServerKey(serverCount int) string {
	serverName := "gogo_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

func (g *gogoAnimeSu) getAnimeListingUrl(serverKey string, pageCount int) string {
	url := server_engine.ParsedServers.GogoanimeServers[serverKey]
	switch serverKey {
	case "gogo_anime_server_1": //https://gogoanime.cm
		return fmt.Sprintf("%s/anime-list.html?page=%d", url, pageCount)
	case "gogo_anime_server_2": //https://gogo-anime.su
		return fmt.Sprintf("%s/all_anime/?page=%d&alphabet=all", url, pageCount)
	case "gogo_anime_server_3": //https://ww2.gogoanimes.org
		return fmt.Sprintf("%s/anime-list?page=%d", url, pageCount)
	case "gogo_anime_server_4": //https://ww1.gogoanime2.org
		return fmt.Sprintf("%s/animelist/all/%d", url, pageCount)
	case "gogo_anime_server_5": //https://www1.gogoanime.sx
		return fmt.Sprintf("%s/anime-list?page=%d", url, pageCount)
	case "gogo_anime_server_6": //https://gogoanime.mom
		return fmt.Sprintf("%s/anime-list", url)
	case "gogo_anime_server_7": //https://gogoanime.wiki
		return fmt.Sprintf("%s/anime-list.html?page=%d", url, pageCount)
	default:
		return ""
	}
}

func (g *gogoAnimeSu) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := g.getAnimeListingUrl(serverKey, pageCount)
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer`),
		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_list_body ul.listing`, &innerHTML),
	}
	//run the assigned task
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

func (g *gogoAnimeSu) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	innerHTML := value["innerHTML"].(string)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(err.Error()), nil
	}

	animeListings := make([]repository.AnimeListing, 0)
	doc.Find(`li`).Each(func(i int, s *goquery.Selection) {
		//get the title from the list
		title := strings.TrimSpace(s.Text())
		navigationUrl := s.Find(`a`).AttrOr("href", "empty")
		fmt.Printf("%s\n%s\n", title, navigationUrl)
		animeListing := repository.AnimeListing{
			NavigationUrl:     navigationUrl,
			AnimeDisplayImage: "",
			AnimeTitle:        title,
		}

		animeListings = append(animeListings, animeListing)
	})

	fmt.Println(value)
	return nil, animeListings
}

//get the anime info
func (g *gogoAnimeSu) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.AnimeInfoSelector(serverCount, path)
}

//anime info html parser
func (g *gogoAnimeSu) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	return GetGogoAnimeServerCm.AnimeInfoHtmlParser(value)
}

//get the lsit of episodes
func (g *gogoAnimeSu) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.EpisodesSelector(serverCount, pageCount, path)
}

//parse the episodes list
func (g *gogoAnimeSu) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.EpisodesHtmlParser(value)
}

//get the episodes info
func (g *gogoAnimeSu) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.EpisodesInfoSelector(serverCount, path)
}

//parse the episodes info
func (g *gogoAnimeSu) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	return GetGogoAnimeServerCm.EpisodesInfoHtmlParser(value)
}

func (g *gogoAnimeSu) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (g *gogoAnimeSu) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}

/*
anime listing server tasks
*/
// type gogoAnimeListingServerTasks struct {
// 	InnerHTML   string
// 	InnerHTMLLi string
// }

// func (gs *gogoAnimeListingServerTasks) gogoanimecm(url string) chromedp.Tasks {
// 	return chromedp.Tasks{
// 		chromedp.Navigate(url),
// 		chromedp.WaitVisible(`footer`),
// 		chromedp.Evaluate(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_list_body ul.listing li').forEach((event) => {event.dispatchEvent(new Event('mouseover'))})`, nil),
// 		chromedp.WaitVisible(`div.tooltip`),
// 		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_list_body ul.listing`, &gs.InnerHTML),
// 		chromedp.InnerHTML(`body`, &gs.InnerHTMLLi, chromedp.ByQuery),
// 	}
// }

// func (gs *gogoAnimeListingServerTasks) gogoanimesu(url string) chromedp.Tasks {
// 	return chromedp.Tasks{
// 		chromedp.Navigate(url),
// 		chromedp.WaitVisible(`footer`),
// 		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_list_body ul.listing`, &gs.InnerHTML),
// 	}
// }

// func (gs *gogoAnimeListingServerTasks) gogoanimeorg(url string) chromedp.Tasks {
// 	return gs.gogoanimecm(url)
// }

// func (gs *gogoAnimeListingServerTasks) gogoanime2org(url string) chromedp.Tasks {
// 	return gs.gogoanimesu(url)
// }

// func (gs *gogoAnimeListingServerTasks) gogoanimesx(url string) chromedp.Tasks {
// 	return chromedp.Tasks{
// 		chromedp.Navigate(url),
// 		chromedp.WaitVisible(`footer`),
// 		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_list_body ul.listing`, &gs.InnerHTML),
// 	}
// }

// func (gs *gogoAnimeListingServerTasks) gogoanimemom(url string, page int) chromedp.Tasks {
// 	js := fmt.Sprintf(`document.querySelector('section.content section.content_left div.main_body.box-anime-list div.anime_name.anime_list ul.pagination-list li a[data-page="%d"]').click()`, page)
// 	return chromedp.Tasks{
// 		chromedp.Navigate(url),
// 		chromedp.WaitVisible(`footer`),
// 		chromedp.Evaluate(js, nil),
// 		chromedp.Sleep(2 * time.Second),
// 		chromedp.InnerHTML(`section.content section.content_left div.main_body.box-anime-list div.anime_list_body ul.listing`, &gs.InnerHTML),
// 	}
// }

// func (gs *gogoAnimeListingServerTasks) gogoanimewiki(url string) chromedp.Tasks {
// 	return gs.gogoanimesx(url)
// }
