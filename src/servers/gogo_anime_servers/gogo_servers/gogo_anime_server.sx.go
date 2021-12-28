package gogo_servers

/*
This source file contains the scraping of gogoanimeserver.sx
*/

import (
	// "fmt"

	// "github.com/chromedp/chromedp"
	"bytes"
	"context"
	"fmt"
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
	GetGogoAnimeServerSx server_engine.AnimeServerInterface = &gogoAnimeSx{}
)

type gogoAnimeSx struct{}

func (g *gogoAnimeSx) getServerKey(serverCount int) string {
	serverName := "gogo_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

func (g *gogoAnimeSx) getAnimeListingUrl(serverKey string, pageCount int) string {
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

func (g *gogoAnimeSx) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
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

func (g *gogoAnimeSx) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
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

	return nil, animeListings
}

//get the anime info
func (g *gogoAnimeSx) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], path)

	var innerHtml string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.InnerHTML(`section.content section.content_left div.main_body#watch div.anime_video_body`, &innerHtml),
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
func (g *gogoAnimeSx) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	innerHTML := value["innerHTML"].(string)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(err.Error()), nil
	}

	imageUrl := ""
	title := doc.Find(`h1`).Text()
	animeType := ""
	summary := ""
	genre := make([]string, 0)
	released := ""
	doc.Find(`div.anime_video_body_cate div.anime_info dl`).Each(func(i int, s *goquery.Selection) {
		spanType := strings.Split(strings.TrimSpace(s.Find(`dt`).Text()), ":")
		spanValue := strings.Split(strings.TrimSpace(s.Find(`dd`).Text()), " ")
		fmt.Println(spanType)
		fmt.Println(strings.TrimSpace(s.Find(`dd`).Text()))
		//split type
		if spanType[0] == "Type" {
			animeType = spanValue[0]
		}
		//split status
		if spanType[0] == "Status" {
			released = spanValue[0]
		}
	})

	//genre
	doc.Find(`div.anime_video_body_cate div.anime_info dl dt`).Each(func(i int, s *goquery.Selection) {
		if s.Text() == "Genre:" {
			doc.Find(`div.anime_video_body_cate div.anime_info dl dd a`).Each(func(i int, s *goquery.Selection) {
				value := s.AttrOr("href", "")
				valueList := strings.Split(value, "/")
				if valueList[len(valueList)-2] == "genre" {
					genre = append(genre, valueList[len(valueList)-1])
				}
			})
		}
	})

	//summary
	summary = doc.Find(`div.anime_video_body_cate div.anime_info div.desc`).Text()

	fmt.Printf("%s\n%s\n%s\n%s\n%s\n%s\n", imageUrl, title, animeType, summary, genre, released)

	animeInfo := repository.AnimeInfo{
		Title:   title,
		Poster:  imageUrl,
		Type:    animeType,
		Summary: summary,
		Generes: strings.Join(genre, ", "),
		Status:  released,
	}

	return nil, &animeInfo
}

//get the lsit of episodes
func (g *gogoAnimeSx) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], path)
	listSel := fmt.Sprintf(`div#episodes ul[data-range="%d"]`, pageCount)

	var noOfPages int
	var innerHtml string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		//get the no of pages
		chromedp.Evaluate(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_video_body ul#episode_page li a').length`, &noOfPages),
		chromedp.Sleep(2 * time.Second),
		chromedp.InnerHTML(listSel, &innerHtml),
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
func (g *gogoAnimeSx) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	pageCount := value["pageCount"].(int)
	noOfPages := value["noOfPages"].(int)
	innerHTML := value["innerHTML"].(string)
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
		return utility.NewInternalServerError(err.Error()), nil
	}

	episodes := make([]repository.AnimeEpisode, 0)
	doc.Find(`li`).Each(func(i int, s *goquery.Selection) {
		navigationUrl := s.Find(`a`).AttrOr("href", "")
		episodeName := fmt.Sprintf("%s %s", "Episode", s.Find(`a`).AttrOr("data-name-normalized", "empty"))
		dataSrc := ""

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
func (g *gogoAnimeSx) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := g.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.GogoanimeServers[serverKey], path)

	var dataSources []string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(c context.Context) error {
			lengthJs := `document.querySelectorAll('section.content section.content_left div.main_body#watch div.anime_video_body div.anime_muti_link ul li a').length`
			//get the server count
			var serverCount int
			if err := chromedp.Evaluate(lengthJs, &serverCount).Do(c); err != nil {
				return err
			}

			//loop the sel
			for i := 0; i < serverCount; i++ {
				//click
				clickJs := fmt.Sprintf(`document.querySelectorAll('section.content section.content_left div.main_body#watch div.anime_video_body div.anime_muti_link ul li a')[%d].click()`, i)
				if err := chromedp.Evaluate(clickJs, nil).Do(c); err != nil {
					return err
				}
				chromedp.Sleep(2 * time.Second).Do(c)
				var dataSrc string
				if err := chromedp.AttributeValue(`div#load_anime div#player iframe`, "src", &dataSrc, nil).Do(c); err != nil {
					return err
				}
				//append
				dataSources = append(dataSources, dataSrc)
			}
			return nil
		}),
	}

	//run the tasks
	servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	fmt.Println(dataSources)
	return nil, map[string]interface{}{
		"dataSources": dataSources,
	}
}

//parse the episodes info
func (g *gogoAnimeSx) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	return nil, value["dataSources"].([]string)
}

func (g *gogoAnimeSx) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (g *gogoAnimeSx) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}

// func (g *gogoAnimeServerSX) AnimeInfo() {
// 	url := "https://www1.gogoanime.sx/anime/gintama-5kq"

// 	//get the doc

// }

// func (g *gogoAnimeServerSX) Episodes() {

// }

// func (g *gogoAnimeServerSX) EpisodeInfo() {
// 	url := "https://www1.gogoanime.sx/anime/gintama-5kq/ep-12"

// 	var dataSources []string
// 	tasks := chromedp.Tasks{
// 		chromedp.Navigate(url),
// 		chromedp.ActionFunc(func(c context.Context) error {
// 			lengthJs := `document.querySelectorAll('section.content section.content_left div.main_body#watch div.anime_video_body div.anime_muti_link ul li a').length`
// 			//get the server count
// 			var serverCount int
// 			if err := chromedp.Evaluate(lengthJs, &serverCount).Do(c); err != nil {
// 				return err
// 			}

// 			//loop the sel
// 			for i := 0; i < serverCount; i++ {
// 				//click
// 				clickJs := fmt.Sprintf(`document.querySelectorAll('section.content section.content_left div.main_body#watch div.anime_video_body div.anime_muti_link ul li a')[%d].click()`, i)
// 				if err := chromedp.Evaluate(clickJs, nil).Do(c); err != nil {
// 					return err
// 				}
// 				chromedp.Sleep(2 * time.Second).Do(c)
// 				var dataSrc string
// 				if err := chromedp.AttributeValue(`div#load_anime div#player iframe`, "src", &dataSrc, nil).Do(c); err != nil {
// 					return err
// 				}
// 				//append
// 				dataSources = append(dataSources, dataSrc)
// 			}
// 			return nil
// 		}),
// 	}

// 	//run the tasks
// 	cancel := g.intializeChromeDpContext(func() chromedp.Tasks {
// 		return tasks
// 	})
// 	defer cancel()

// 	fmt.Println(dataSources)

// 	// //get the doc
// 	// byteData := []byte(innerHtml)
// 	// readerData := bytes.NewReader(byteData)

// 	// //process the html and parse it to the main json value
// 	// doc, err := goquery.NewDocumentFromReader(readerData)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// doc.Find(``)
// }
