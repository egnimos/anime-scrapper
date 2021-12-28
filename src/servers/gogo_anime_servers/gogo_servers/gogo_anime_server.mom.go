package gogo_servers

/*
This source file contains the scraping of gogoanimeserver.mom
*/

import (
	// "fmt"

	// "github.com/chromedp/chromedp"
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/servers"

	// "github.com/egnimos/anime-scrapper/src/servers"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetGogoAnimeServerMom server_engine.AnimeServerInterface = &gogoAnimeMom{}
)

type gogoAnimeMom struct{}

func (g *gogoAnimeMom) getServerKey(serverCount int) string {
	serverName := "gogo_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

func (g *gogoAnimeMom) getAnimeListingUrl(serverKey string, pageCount int) string {
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

func (g *gogoAnimeMom) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	//if the pageCount is 0
	if pageCount == 0 {
		pageCount = 1
	}
	serverKey := g.getServerKey(serverCount)
	url := g.getAnimeListingUrl(serverKey, pageCount)
	var innerHTML string
	js := fmt.Sprintf(`document.querySelector('section.content section.content_left div.main_body.box-anime-list div.anime_name.anime_list ul.pagination-list li a[data-page="%d"]').click()`, pageCount)
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer`),
		chromedp.Evaluate(js, nil),
		chromedp.Sleep(2 * time.Second),
		chromedp.InnerHTML(`section.content section.content_left div.main_body.box-anime-list div.anime_list_body ul.listing`, &innerHTML),
	}

	// run the assigned task
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

func (g *gogoAnimeMom) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	innerHTML := value["innerHTML"].(string)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	animeListings := make([]repository.AnimeListing, 0)
	doc.Find(`li`).Each(func(i int, s *goquery.Selection) {
		//get the title from the list
		title := strings.TrimSpace(s.Text())
		navigationUrl := strings.Split(s.Find(`a`).AttrOr("href", ""), "https://gogoanime.mom")
		fmt.Printf("%s\n%s\n", title, navigationUrl)
		animeListing := repository.AnimeListing{
			NavigationUrl:     navigationUrl[len(navigationUrl)-1],
			AnimeDisplayImage: "",
			AnimeTitle:        title,
		}

		animeListings = append(animeListings, animeListing)
	})

	return nil, animeListings
}

//get the anime info
func (g *gogoAnimeMom) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], path)
	var innerHtml string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer`),
		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_video_body`, &innerHtml),
	}

	//run the tasks
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return nil, map[string]interface{}{
		"innerHTML": innerHtml,
	}
}

//anime info html parser
func (g *gogoAnimeMom) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	innerHTML := value["innerHTML"].(string)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	animeInfo := repository.AnimeInfo{}
	genres := make([]string, 0)
	doc.Find(`div.anime_video_body_cate a`).Each(func(i int, s *goquery.Selection) {
		genre := strings.TrimSpace(s.AttrOr("title", ""))
		genres = append(genres, genre)
	})

	animeInfo.Title = strings.TrimSpace(doc.Find(`div.anime_video_body_cate div.anime-info a`).Text())
	animeInfo.Generes = strings.Join(genres, ", ")

	return nil, &animeInfo
}

//get the list of episodes
func (g *gogoAnimeMom) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], path)
	js := fmt.Sprintf(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_video_body ul#episode_page li a')[%d].click()`, pageCount)

	var noOfPages int
	var innerHtml string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Evaluate(js, nil), //click to get the next set of episodes
		//get the no of pages
		chromedp.Evaluate(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_video_body ul#episode_page li a').length`, &noOfPages),
		chromedp.Sleep(2 * time.Second),
		chromedp.InnerHTML(`div#load_ep ul#episode_related`, &innerHtml),
	}

	//run the tasks
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return nil, map[string]interface{}{
		"innerHTML": innerHtml,
		"noOfPages": noOfPages,
		"pageCount": pageCount,
		"path":      path,
	}
}

//parse the episodes list
func (g *gogoAnimeMom) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	pageCount := value["pageCount"].(int)
	noOfPages := value["noOfPages"].(int)
	innerHTML := value["innerHTML"].(string)
	path := value["path"].(string)
	paginationStatus := "available"
	//check the page cond
	if pageCount+1 == noOfPages {
		paginationStatus = "disabled"
	}

	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		panic(err)
	}

	episodes := make([]repository.AnimeEpisode, 0)
	doc.Find(`li`).Each(func(i int, s *goquery.Selection) {
		navigationUrl := fmt.Sprintf("%s%s", path, s.Find(`a`).AttrOr("data-order", ""))
		episodeName := fmt.Sprintf("%s %s %s", "Episode", s.Find(`a`).AttrOr("data-order", ""), s.Find(`a div.cate`).Text())
		dataSrc := s.Find(`a`).AttrOr("data-src", "")

		fmt.Printf("%s\n%s\n%s\n", navigationUrl, episodeName, dataSrc)
		episode := repository.AnimeEpisode{
			EpisodePath: navigationUrl,
			Episode:     episodeName,
			DataSrc:     dataSrc,
		}

		episodes = append(episodes, episode)
	})

	return nil, map[string]interface{}{
		"pagination_status": paginationStatus,
		"episodes":          episodes,
	}
}

//get the episodes info
func (g *gogoAnimeMom) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	values := strings.Split(path, "/")
	paths := strings.Split(path, values[len(values)-1])
	fmt.Println(values[len(values)-1])
	fmt.Println(paths[0])
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], paths[0])
	var dataSrc string
	value, err := strconv.Atoi(values[len(values)-1])
	if err != nil {
		return utility.NewBadRequestError(err.Error()), nil
	}

	//tasks
	tasks := chromedp.Tasks{
		//navigate
		chromedp.Navigate(url),
		//perform the action
		chromedp.ActionFunc(func(c context.Context) error {
			//get the innerhtml of the given episode page
			var innerHtml string
			if err := chromedp.InnerHTML(`div.anime_video_body ul#episode_page`, &innerHtml).Do(c); err != nil {
				return err
			}

			byteData := []byte(innerHtml)
			readerData := bytes.NewReader(byteData)

			//process the html and parse it to the main json value
			doc, err := goquery.NewDocumentFromReader(readerData)
			if err != nil {
				return err
			}

			doc.Find(`li a`).EachWithBreak(func(i int, s *goquery.Selection) bool {
				js := fmt.Sprintf(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_video_body ul#episode_page li a')[%d].click()`, i)
				rangeValue := strings.Split(strings.TrimSpace(s.Text()), "-")
				// fmt.Print(rangeValue)
				value1, _ := strconv.Atoi(rangeValue[0])
				value2, _ := strconv.Atoi(rangeValue[1])
				//check the range
				if value1 <= value && value <= value2 {
					if err := chromedp.Evaluate(js, nil).Do(c); err != nil {
						return false
					}

					//get the datasrc
					sel := fmt.Sprintf(`div#load_ep ul#episode_related li a[data-order="%d"]`, value)
					if err := chromedp.AttributeValue(sel, "data-src", &dataSrc, nil).Do(c); err != nil {
						return false
					}

					return false

				}
				return true
			})

			return nil
		}),
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
func (g *gogoAnimeMom) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	dataSrc := value["dataSource"].(string)
	dataSources := []string{dataSrc}
	return nil, dataSources
}

func (g *gogoAnimeMom) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (g *gogoAnimeMom) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}

// func (g *gogoAnimeServer6) AnimeInfo() {
// 	url := "https://gogoanime.mom/movie/gintama/"

// 	//get the doc

// }

// func (g *gogoAnimeServer6) Episodes() {
// 	// var pageCount int

// }

// func (g *gogoAnimeServer6) EpisodeInfo() {
// 	url := "https://gogoanime.mom/movie/gintama/"

// 	// //get the doc
// 	// byteData := []byte(innerHtml)
// 	// readerData := bytes.NewReader(byteData)

// 	// //process the html and parse it to the main json value
// 	// doc, err := goquery.NewDocumentFromReader(readerData)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// dataSrc := doc.Find(`li a[data-order="100"]`).AttrOr("data-src", "")
// 	// fmt.Printf("%s\n", dataSrc)
// }
