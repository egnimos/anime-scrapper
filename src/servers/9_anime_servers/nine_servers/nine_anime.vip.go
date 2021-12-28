package nine_servers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/servers"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetNineAnimeVIP server_engine.AnimeServerInterface = &nineAnimeVIP{}
)

type nineAnimeVIP struct{}

func (n *nineAnimeVIP) getServerKey(serverCount int) string {
	serverName := "9_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

//get anime listing selectors
func (n *nineAnimeVIP) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	if pageCount == 0 {
		pageCount = 1
	}

	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s/az-list?page=%d", server_engine.ParsedServers.NineanimeServers[serverKey], pageCount)
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div#footer`),
		chromedp.InnerHTML(`div#main div.content div.widget.az-list div.widget-body div.items`, &innerHTML),
	}

	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

func (n *nineAnimeVIP) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	innerHTML := value["innerHTML"].(string)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(err.Error()), nil
	}

	animeListings := make([]repository.AnimeListing, 0)
	doc.Find(`div`).Each(func(i int, s *goquery.Selection) {
		//get the title from the list
		title := strings.TrimSpace(s.Find(`div.info a.name`).Text())
		navigationUrl := s.Find(`div.info a.name`).AttrOr("href", "")
		poster := s.Find(`img.thumb.tooltipstered`).AttrOr("src", "")
		fmt.Printf("%s\n%s\n%s\n", title, navigationUrl, poster)
		animeListing := repository.AnimeListing{
			NavigationUrl:     navigationUrl,
			AnimeDisplayImage: poster,
			AnimeTitle:        title,
		}

		animeListings = append(animeListings, animeListing)
	})

	return nil, animeListings
}

func (n *nineAnimeVIP) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div#footer`),
		chromedp.InnerHTML(`div#main div.content div.widget.info div.widget-body div.row`, &innerHTML),
	}

	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

func (n *nineAnimeVIP) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	innerHTML := value["innerHTML"].(string)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(err.Error()), nil
	}

	var animeInfo repository.AnimeInfo
	animeInfo.Title = doc.Find(`div.info.col-md-19 div.head div.c1 h2.title`).AttrOr("data-jtitle", "")
	animeInfo.Summary = doc.Find(`div.info.col-md-19 div.desc div.long`).Text()
	animeInfo.Poster = doc.Find(`div.thumb.col-md-5.hidden-sm.hidden-xs img`).AttrOr("src", "")
	genres := make([]string, 0)
	doc.Find(`div#main div.info.col-md-19 div.row dl.meta.col-sm-12 dd a`).Each(func(i int, s *goquery.Selection) {
		genres = append(genres, s.AttrOr("title", ""))
	})
	animeInfo.Generes = strings.Join(genres, ", ")
	return nil, &animeInfo
}

func (n *nineAnimeVIP) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)
	episodeList := fmt.Sprintf(`div#servers-container ul.episodes.range.sp[data-range-id="%d"]`, pageCount)
	
	var noOfPages int
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div#footer`),
		chromedp.Evaluate(`document.querySelectorAll('div#servers-container div.range span').length`, &noOfPages),
		chromedp.InnerHTML(episodeList, &innerHTML),
	}

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

func (n *nineAnimeVIP) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
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

func (n *nineAnimeVIP) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div#footer`),
		chromedp.InnerHTML(`div#main div.content div.widget.info div.widget-body div.row`, &innerHTML),
	}

	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

func (n *nineAnimeVIP) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	return nil, nil
}

func (n *nineAnimeVIP) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeVIP) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}
