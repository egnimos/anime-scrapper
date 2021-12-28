package nine_servers

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/egnimos/anime-scrapper/src/servers"
	"github.com/egnimos/anime-scrapper/src/servers/gogo_anime_servers/gogo_servers"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetNineAnimeSurf server_engine.AnimeServerInterface = &nineAnimeSurf{}
)

type nineAnimeSurf struct{}

func (n *nineAnimeSurf) getServerKey(serverCount int) string {
	serverName := "9_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

//get the list of anime
func (n *nineAnimeSurf) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	if pageCount == 0 {
		pageCount = 1
	}

	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s/all_anime/?page=%d&alphabet=all", server_engine.ParsedServers.NineanimeServers[serverKey], pageCount)

	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer`),
		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_list_body ul.listing`, &innerHTML),
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
func (n *nineAnimeSurf) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	return gogo_servers.GetGogoAnimeServerSu.AnimeListingHtmlParser(value)
}

//get the anime info
func (n *nineAnimeSurf) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)

	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_info_body`, &innerHTML),
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
func (n *nineAnimeSurf) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	return gogo_servers.GetGogoAnimeServerSu.AnimeInfoHtmlParser(value)
}

//get the lsit of episodes
func (n *nineAnimeSurf) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)
	// js := fmt.Sprintf(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_video_body ul#episode_page li a')[%d].click()`, pageCount)

	var noOfPages int
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		//click to get the next set of episodes
		// chromedp.Evaluate(js, nil),
		//get the no of pages
		chromedp.Evaluate(`document.querySelectorAll('section.content section.content_left div.main_body div.anime_video_body ul#episode_page li a').length`, &noOfPages),
		chromedp.InnerHTML(`div#load_ep ul#episode_related`, &innerHTML),
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
func (n *nineAnimeSurf) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	return gogo_servers.GetGogoAnimeServerSu.EpisodesHtmlParser(value)
}

//get the episodes info
func (n *nineAnimeSurf) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	serverKey := n.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.NineanimeServers[serverKey], path)
	var innerHtml string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.InnerHTML(`section.content section.content_left div.main_body div.anime_video_body div.anime_muti_link`, &innerHtml),
	}

	//run the tasks
	if err := servers.InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	}); err != nil {
		return err, nil
	}

	return nil, map[string]interface{}{
		"innerHTML": innerHtml,
	}
}

//parse the episodes info
func (n *nineAnimeSurf) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	return gogo_servers.GetGogoAnimeServerSu.EpisodesInfoHtmlParser(value)
}

func (n *nineAnimeSurf) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (n *nineAnimeSurf) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}
