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
	GetNineAnime2Com server_engine.AnimeServerInterface = &nineAnime2Com{}
)

type nineAnime2Com struct{}

func (n *nineAnime2Com) getServerKey(serverCount int) string {
	serverName := "9_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

//get the list of anime
func (n *nineAnime2Com) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s/home/%d", server_engine.ParsedServers.NineanimeServers[serverKey], pageCount)

	var innerHTML string
	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.WaitVisible(`footer`),
			chromedp.InnerHTML(`section.nobg div.body ul.anime-list`, &innerHTML),
		}
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

//parse the list of anime
func (n *nineAnime2Com) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
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
		displayImage := s.Find(`a.poster.tooltipstered img`).AttrOr("src", "")
		title := s.Find(`a.name`).AttrOr("data-jtitle", "")
		navigationUrl := s.Find(`a.name`).AttrOr("href", "")

		animeListing := repository.AnimeListing{
			NavigationUrl:     navigationUrl,
			AnimeDisplayImage: displayImage,
			AnimeTitle:        title,
		}

		animeListings = append(animeListings, animeListing)
	})

	return nil, animeListings
}

//get the anime info
func (n *nineAnime2Com) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)

	var innerHTML string
	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.WaitVisible(`footer`),
			chromedp.InnerHTML(`section#info`, &innerHTML),
		}
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"innerHTML":   innerHTML,
		"serverCount": serverCount,
	}
}

//anime info html parser
func (n *nineAnime2Com) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	innerHTML := value["innerHTML"].(string)
	serverCount := value["serverCount"].(int)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	animeInfo := &repository.AnimeInfo{}

	//title
	animeInfo.Title = doc.Find(`div.info h1.title`).AttrOr("data-jtitle", "")
	//poster
	serverKey := n.getServerKey(serverCount)
	animeInfo.Poster = fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], doc.Find(`div.thumb div img`).AttrOr("src", ""))
	animeInfo.Summary = doc.Find(`div.info p.shorting`).Text()

	doc.Find(`div.info div.meta div.col1 div`).Each(func(i int, s *goquery.Selection) {
		//text
		switchType := strings.TrimSpace(s.Text())
		switch switchType {
		case "Type:":
			animeInfo.Type = s.Find(`span a`).Text()
		case "Status:":
			animeInfo.Status = s.Find(`span`).Text()
		case "Genre:":
			{
				genres := make([]string, 0)
				s.Find(`span a`).Each(func(i int, s *goquery.Selection) {
					genres = append(genres, s.AttrOr("title", ""))
				})
				animeInfo.Generes = strings.Join(genres, ",")
			}
		default:
		}
	})

	return nil, animeInfo
}

//get the lsit of episodes
func (n *nineAnime2Com) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)

	var innerHTML string
	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.WaitVisible(`footer`),
			chromedp.InnerHTML(`div#episodes section div.body`, &innerHTML),
		}
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

//parse the episodes list
func (n *nineAnime2Com) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	innerHTML := value["innerHTML"].(string)
	byteData := []byte(innerHTML)
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	episodes := make([]repository.AnimeEpisode, 0)
	doc.Find(`ul.episodes li`).Each(func(i int, s *goquery.Selection) {
		episode := repository.AnimeEpisode{}
		episode.EpisodePath = s.Find(`a`).AttrOr("href", "")
		episode.Episode = fmt.Sprintf("%s %s", "Episode", s.Find(`a`).Text())

		episodes = append(episodes, episode)
	})

	return nil, map[string]interface{}{
		"pagination_status": "disabled",
		"episodes":          episodes,
	}
}

//get the episodes info
func (n *nineAnime2Com) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)

	var dataSource string
	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.WaitVisible(`footer`),
			chromedp.AttributeValue(`div#player iframe#playerframe`, "src", &dataSource, nil),
		}
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"dataSource": dataSource,
	}
}

//parse the episodes info
func (n *nineAnime2Com) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	dataSources := []string{value["dataSource"].(string)}
	return nil, dataSources
}

func (n *nineAnime2Com) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnime2Com) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}
