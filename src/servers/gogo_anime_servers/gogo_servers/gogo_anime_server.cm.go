package gogo_servers

/*
This source file contains the scraping of gogoanimeserver.cm
*/

import (
	// "fmt"

	// "github.com/chromedp/chromedp"
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
	GetGogoAnimeServerCm server_engine.AnimeServerInterface = &gogoAnimeCm{}
)

type gogoAnimeCm struct{}

func (g *gogoAnimeCm) getServerKey(serverCount int) string {
	serverName := "gogo_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

func (g *gogoAnimeCm) getAnimeListingUrl(serverKey string, pageCount int) string {
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

func (g *gogoAnimeCm) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := g.getAnimeListingUrl(serverKey, pageCount)
	var innerHTML string
	var innerHTMLLi string

	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer`),
		chromedp.Evaluate(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_list_body ul.listing li').forEach((event) => {event.dispatchEvent(new Event('mouseover'))})`, nil),
		chromedp.WaitVisible(`div.tooltip`),
		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_list_body ul.listing`, &innerHTML),
		chromedp.InnerHTML(`body`, &innerHTMLLi, chromedp.ByQuery),
	}
	// //run the assigned task
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return nil, map[string]interface{}{
		"innerHTML":   innerHTML,
		"innerHTMLLi": innerHTMLLi,
	}
}

func (g *gogoAnimeCm) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	innerHTML := value["innerHTML"].(string)
	innerHTMLLi := value["innerHTMLLi"].(string)
	//get the doc
	byteData := []byte(innerHTML)
	byteData1 := []byte(innerHTMLLi)
	readerData := bytes.NewReader(byteData)
	readerData1 := bytes.NewReader(byteData1)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}
	//process for reader1
	doc1, err := goquery.NewDocumentFromReader(readerData1)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	animeListings := make([]repository.AnimeListing, 0)
	doc1.Find(`div.tooltip`).Each(func(i int, ts *goquery.Selection) {
		doc.Find(`li`).EachWithBreak(func(i int, s *goquery.Selection) bool {
			//get the title from the list
			mainTitle := strings.TrimSpace(s.Text())
			title := ts.Find(`div a.bigChar`).Text()
			//check whether the keyword is contains or not
			if mainTitle == title {
				imageUrl := ts.Find(`div.thumnail_tool img`).AttrOr("src", "empty")
				navigationUrl := s.Find(`a`).AttrOr("href", "empty")
				fmt.Printf("%s\n%s\n%s\n%s\n", imageUrl, title, mainTitle, navigationUrl)
				animeListing := repository.AnimeListing{
					NavigationUrl:     navigationUrl,
					AnimeDisplayImage: imageUrl,
					AnimeTitle:        title,
				}

				animeListings = append(animeListings, animeListing)
				return false
			}
			return true
		})
	})

	return nil, animeListings
}

//get the anime info
func (g *gogoAnimeCm) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], path)

	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_info_body`, &innerHTML),
	}

	//run the assigned task
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

//anime info html parser
func (g *gogoAnimeCm) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	innerHTML := value["innerHTML"].(string)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	var animeInfo repository.AnimeInfo

	animeInfo.Poster = doc.Find(`div.anime_info_body_bg img`).AttrOr("src", "empty")
	animeInfo.Title = doc.Find(`div.anime_info_body_bg h1`).Text()
	genre := make([]string, 0)
	doc.Find(`div.anime_info_body_bg p.type`).Each(func(i int, s *goquery.Selection) {
		//get the span
		spanType := strings.TrimSpace(s.Find(`span`).Text())
		switch spanType {
		case "Type:":
			{
				animeInfo.Type = s.Find(`a`).Text()
			}
		case "Plot Summary:":
			{
				animeInfo.Summary = s.Text()
			}
		case "Genre:":
			{
				s.Find(`a`).Each(func(i int, s *goquery.Selection) {
					genre = append(genre, s.Text())
				})
				animeInfo.Generes = strings.Join(genre, ", ")
			}
		case "Released:":
			{
				animeInfo.Status = s.Text()
			}
		}
	})

	return nil, &animeInfo
}

//get the lsit of episodes
func (g *gogoAnimeCm) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], path)
	js := fmt.Sprintf(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_video_body ul#episode_page li a')[%d].click()`, pageCount)

	var noOfPages int
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		//click to get the next set of episodes
		chromedp.Evaluate(js, nil),
		//get the no of pages
		chromedp.Evaluate(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_video_body ul#episode_page li a').length`, &noOfPages),
		chromedp.InnerHTML(`div#load_ep ul#episode_related`, &innerHTML),
	}

	//run the assigned task
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
		"noOfPages": noOfPages,
		"pageCount": pageCount,
	}
}

//parse the episodes list
func (g *gogoAnimeCm) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	innerHTML := value["innerHTML"].(string)
	noOfPages := value["noOfPages"].(int)
	pageCount := value["pageCount"].(int)
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
		"pagination_status": paginationStatus,
		"episodes":          animeEpisodes,
	}
}

//get the episodes info
func (g *gogoAnimeCm) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], path)
	var innerHtml string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_video_body div.anime_muti_link`, &innerHtml),
	}

	//run the tasks
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return nil, map[string]interface{}{
		"innerHTML": innerHtml,
	}
}

//parse the episodes info
func (g *gogoAnimeCm) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	innerHTML := value["innerHTML"].(string)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(err.Error()), nil
	}

	dataSrc := make([]string, 0)
	doc.Find(`ul li a`).Each(func(i int, s *goquery.Selection) {
		dataSrc = append(dataSrc, s.AttrOr("data-video", ""))
	})

	fmt.Println(dataSrc)
	return nil, dataSrc
}

func (g *gogoAnimeCm) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (g *gogoAnimeCm) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}
