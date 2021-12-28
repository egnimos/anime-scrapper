package nine_servers

import (
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
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetNineAnimeCz server_engine.AnimeServerInterface = &nineAnimeCz{}
)

type nineAnimeCz struct{}

func (n *nineAnimeCz) getServerKey(serverCount int) string {
	serverName := "9_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

// https://9anime.cz/az-list?page=2

//get the list of anime
func (n *nineAnimeCz) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	if pageCount == 0 {
		pageCount = 1
	}

	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s/az-list?page=%d", server_engine.ParsedServers.NineanimeServers[serverKey], pageCount)

	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer`),
		chromedp.InnerHTML(`div#body section div.body ul.anime-list-v`, &innerHTML),
	}
	//run the assigned task
	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

//parse the list of anime
func (n *nineAnimeCz) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
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
		title := strings.TrimSpace(s.Find(`div.info a.name`).AttrOr("data-jtitle", ""))
		navigationUrl := s.Find(`div.info a.name`).AttrOr("href", "")
		poster := s.Find(`div.thumb div.tooltipstered img`).AttrOr("src", "")
		fmt.Printf("%s\n%s\n", title, navigationUrl)
		animeListing := repository.AnimeListing{
			NavigationUrl:     navigationUrl,
			AnimeDisplayImage: poster,
			AnimeTitle:        title,
		}

		animeListings = append(animeListings, animeListing)
	})

	return nil, animeListings
}

//get the anime info
func (n *nineAnimeCz) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)

	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer`),
		chromedp.InnerHTML(`section#info`, &innerHTML),
	}

	//run the assigned task
	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

//anime info html parser
func (n *nineAnimeCz) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	innerHTML := value["innerHTML"].(string)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(err.Error()), nil
	}

	var animeInfo repository.AnimeInfo
	animeInfo.Title = doc.Find(`div.info h1.title`).AttrOr("data-jtitle", "")
	animeInfo.Summary = doc.Find(`div.info p.shorting`).Text()
	animeInfo.Poster = doc.Find(`div.thumb div img`).AttrOr("src", "")
	genres := make([]string, 0)
	doc.Find(`div.info div.meta div.col1 div`).Each(func(i int, s *goquery.Selection) {
		typ := strings.TrimSpace(s.Text())
		switch typ {
		case "Type:":
			animeInfo.Type = s.Find(`span a`).Text()
		case "Status:":
			animeInfo.Status = s.Find(`span`).Text()
		case "Genre:":
			{
				s.Find(`span a`).Each(func(i int, s *goquery.Selection) {
					genres = append(genres, s.AttrOr("title", ""))
				})

				animeInfo.Generes = strings.Join(genres, ", ")
			}
		}
	})

	return nil, &animeInfo
}

//get the lsit of episodes
func (n *nineAnimeCz) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)
	episodesInnerHTML := fmt.Sprintf(`div#episodes section div.body ul.episodes[data-range="%d"]`, pageCount)

	var noOfPages int
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer`),
		chromedp.WaitVisible(`div#episodes section div.head`),
		//get the no of pages
		chromedp.Evaluate(`document.querySelectorAll('div#episodes section div.head ul.ranges li').length`, &noOfPages),
		chromedp.InnerHTML(episodesInnerHTML, &innerHTML),
	}

	//run the assigned task
	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
		"noOfPages": noOfPages,
		"pageCount": pageCount,
	}
}

//parse the episodes list
func (n *nineAnimeCz) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	pageCount := value["pageCount"].(int)
	noOfPages := value["noOfPages"].(int)
	innerHTML := value["innerHTML"].(string)
	paginationStatus := "available"
	//check the page cond
	if pageCount+1 >= noOfPages {
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
		episodeName := fmt.Sprintf("%s %s", "Episode", s.Find(`a`).Text())

		fmt.Printf("%s\n%s\n", navigationUrl, episodeName)
		episode := repository.AnimeEpisode{
			EpisodePath: navigationUrl,
			Episode:     episodeName,
		}

		episodes = append(episodes, episode)
	})

	return nil, map[string]interface{}{
		"pagination_status": paginationStatus,
		"episodes":          episodes,
	}
}

//get the episodes info
func (n *nineAnimeCz) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)

	var dataSources []string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer`),
		chromedp.WaitVisible(`div#episodes section div.head`),
		chromedp.ActionFunc(func(c context.Context) error {
			serverLength := `document.querySelectorAll('#episodes section div.tabs.servers.notab span').length`
			//get the server count
			var serverCount int
			if err := chromedp.Evaluate(serverLength, &serverCount).Do(c); err != nil {
				return err
			}

			//loop the sel
			for i := 0; i < serverCount; i++ {
				//click
				selectServer := fmt.Sprintf(`document.querySelectorAll('#episodes section div.tabs.servers.notab span')[%d].click()`, i)
				if err := chromedp.Evaluate(selectServer, nil).Do(c); err != nil {
					return err
				}
				chromedp.Sleep(2 * time.Second).Do(c)
				var dataSrc string
				if err := chromedp.AttributeValue(`div#player iframe`, "src", &dataSrc, nil).Do(c); err != nil {
					return err
				}
				//append
				dataSources = append(dataSources, dataSrc)
			}
			return nil
		}),
	}

	//run the tasks
	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	}); err != nil {
		return err, nil
	}

	fmt.Println(dataSources)
	return nil, map[string]interface{}{
		"dataSources": dataSources,
	}
}

//parse the episodes info
func (n *nineAnimeCz) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	return nil, value["dataSources"].([]string)
}

func (n *nineAnimeCz) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeCz) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}
