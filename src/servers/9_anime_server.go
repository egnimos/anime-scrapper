package servers

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetNineAnimeServers server_engine.AnimeServerInterface = &nineAnime{}
)

type nineAnime struct {
}

func (n *nineAnime) getServerKey(serverCount int) string {
	serverName := "9_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

func (n *nineAnime) AnimeListingSelector(serverCount int, pageCount int) map[string]interface{} {
	serverKey := n.getServerKey(serverCount)
	url := server_engine.ParsedServers.NineanimeServers[serverKey]
	var innerHTML string
	InitializeChromeDp(func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.WaitVisible(`footer#footer`),
			chromedp.InnerHTML(`#main-video-list div.video-list.row.mx-0`, &innerHTML),
		}
	})

	return map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

func (n *nineAnime) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	return nil, nil
}

//get the anime info
func (n *nineAnime) AnimeInfoSelector(serverCount int, path string) map[string]interface{} {
	return nil
}

//anime info html parser
func (n *nineAnime) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	return nil, nil
}

//get the lsit of episodes
func (n *nineAnime) EpisodesSelector(serverCount int, pageCount int, path string) map[string]interface{} {
	return nil
}

//parse the episodes list
func (n *nineAnime) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

//get the episodes info
func (n *nineAnime) EpisodesInfoSelector(serverCount int, path string) map[string]interface{} {
	return nil
}

//parse the episodes info
func (n *nineAnime) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	return nil, nil
}
