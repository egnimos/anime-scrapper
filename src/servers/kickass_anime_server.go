package servers

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
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetKickAssServer server_engine.AnimeServerInterface = &kickass{}
)

type kickass struct{}

func (k *kickass) getServerKey(serverCount int) string {
	serverName := "kickass_anime_server"
	return fmt.Sprintf("%s_%d", serverName, serverCount)
}

//get the list of anime
func (k *kickass) AnimeListingSelector(serverCount int, pageCount int) map[string]interface{} {
	server := k.getServerKey(serverCount)
	url := fmt.Sprintf("%s/anime-list", server_engine.ParsedServers.KickassanimeServers[server])
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer#footer`),
		//display it to the grid format
		chromedp.Click(`input#__BVID__50__BV_option_1_`, chromedp.ByQuery),
		//display 100 items
		//*[@id="__BVID__47"]//option[contains(@value, "100")]
		// chromedp.Click(`//select[@id="__BVID__47"]/option[@value="100"]`, chromedp.BySearch),
		//wait for 2 sec
		// chromedp.Sleep(2 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			clickCount := pageCount
			for i := 0; i < clickCount; i++ {
				if err := chromedp.Run(ctx,
					//click the next
					chromedp.Click(`#content button.btn.btn-primary svg.svg-inline--fa.fa-chevron-right.fa-w-10`, chromedp.ByQueryAll),
					//sleep
					chromedp.Sleep(2*time.Second)); err != nil {
					return err
				}
			}

			return nil
		}),
		chromedp.InnerHTML(`#content div.row.video-list.row.mx-0 div.col.row.video-list.row.mx-0`, &innerHTML),
	}
	//run the assigned task
	InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

//parse anime listing html return value
func (k *kickass) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
	innerHTML := value["innerHTML"]
	byteData := []byte(innerHTML.(string))
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	animeListings := make([]repository.AnimeListing, 0)
	doc.Find("div.video-item.col-6.mb-2.px-1.col-2").Each(func(i int, s *goquery.Selection) {
		navigationUrl := s.Find("a.ka-url-wrapper.video-item-poster.rounded").AttrOr("href", "empty")
		bgUrl := s.Find("a.ka-url-wrapper.video-item-poster.rounded").AttrOr("style", "empty")
		imageUrl := bytes.Split([]byte(bgUrl), []byte("\""))
		title := s.Find("a.ka-url-wrapper.video-item-title").Text()
		// fmt.Printf("%s\n%s\n%s\n", navigationUrl, imageUrl[1], title)
		animeListing := repository.AnimeListing{
			NavigationUrl:     navigationUrl,
			AnimeDisplayImage: string(imageUrl[1]),
			AnimeTitle:        title,
		}

		animeListings = append(animeListings, animeListing)
	})

	return nil, animeListings
}

//get the anime info
func (k *kickass) AnimeInfoSelector(serverCount int, path string) map[string]interface{} {
	server := k.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.KickassanimeServers[server], path)
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer#footer`),
		chromedp.InnerHTML(`div#main.container-fluid.pt-3`, &innerHTML),
	}
	//run the assigned task
	InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

//anime info html parser
func (k *kickass) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	innerHTML := value["innerHTML"]
	byteData := []byte(innerHTML.(string))
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	var animeInfo repository.AnimeInfo
	//1st block
	doc.Find("div#sidebar-anime-info.border.rounded.mb-3").Each(func(i int, s *goquery.Selection) {
		animeImage := s.Find(`div.poster`).AttrOr("style", "empty")
		imageUrl := bytes.Split([]byte(animeImage), []byte("\""))
		animeMap := map[string]interface{}{}
		s.Find(`div.p-3 div.mb-2`).Each(func(i int, s *goquery.Selection) {
			key := s.Find(`div.font-weight-bold.mr-1`).Text()
			value := s.Find(`div.mb-2 span`).Text()
			animeMap[key] = value
		})
		fmt.Printf("%s\n%s\n", imageUrl[1], animeMap)
		animeInfo.Poster = string(imageUrl[1])
		animeInfo.Type = animeMap["Type:"].(string)
		animeInfo.Status = animeMap["Status:"].(string)
	})

	//2nd block
	doc.Find("#content div.anime-info.border.rounded.mb-3").Each(func(i int, s *goquery.Selection) {
		//a.ka-url-wrapper div.info-header.p-3.px-4.hep div.info-wrapper h1.title
		title := s.Find(`a.ka-url-wrapper div.info-header.p-3.px-4.hep div.info-wrapper h1.title`).Text()
		//div.container.p-3 div.mb-3 p.mb-0
		summary := s.Find(`div.container.p-3 div.mb-3 p.mb-0`).Text()
		//div.mb-3 a.ka-url-wrapper.d-inline-block
		genres := make([]string, 0)
		s.Find(`div.mb-3 a.ka-url-wrapper.d-inline-block`).Each(func(i int, s *goquery.Selection) {
			genres = append(genres, s.Text())
		})

		fmt.Printf("%s\n%s\n", summary, genres)
		animeInfo.Generes = strings.Join(genres, ", ")
		animeInfo.Summary = summary
		animeInfo.Title = title
	})

	return nil, &animeInfo
}

//get the list of episodes of the given anime
func (k *kickass) EpisodesSelector(serverCount int, pageCount int, path string) map[string]interface{} {
	server := k.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.KickassanimeServers[server], path)
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer#footer`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			clickCount := pageCount
			for i := 0; i < clickCount; i++ {
				if err := chromedp.Run(ctx,
					//click the next
					chromedp.Click(`#content div.main-episode-list.border.rounded.p-3.mb-3 button.btn.btn-primary`, chromedp.ByQueryAll),
					//sleep
					chromedp.Sleep(2*time.Second)); err != nil {
					return err
				}
			}

			return nil
		}),
		chromedp.Click(`#content div.main-episode-list.border.rounded.p-3.mb-3 table.table.b-table.table-hover thead[role="rowgroup"] tr[role="row"] th[aria-colindex="1"] span.sr-only`, chromedp.ByQuery),
		chromedp.InnerHTML(`#content div.main-episode-list.border.rounded.p-3.mb-3`, &innerHTML),
	}
	//run the assigned task
	InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

//parse the episodes html
func (k *kickass) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	innerHTML := value["innerHTML"]
	byteData := []byte(innerHTML.(string))
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	//parse the html info and convert it to the json file
	animeEpisodes := make([]repository.AnimeEpisode, 0)
	doc.Find(`table.table.b-table.table-hover tbody[role="rowgroup"] tr[role="row"]`).Each(func(i int, s *goquery.Selection) {
		//episode path
		episodesPathUrl := s.Find(`td[aria-colindex="1"] a.ka-url-wrapper`).AttrOr("href", "empty")
		episode := s.Find(`td[aria-colindex="1"] a.ka-url-wrapper`).Text()
		fmt.Printf("%s\n%s\n", episodesPathUrl, episode)
		animeEpisode := repository.AnimeEpisode{
			EpisodePath: episodesPathUrl,
			Episode:     episode,
		}

		animeEpisodes = append(animeEpisodes, animeEpisode)
	})

	//check if the next button is enabled then we can paginate
	paginateStatus := "not-disabled"
	doc.Find(`div.row.align-items-end div.text-right.col-lg-3.col-12 div.mb-3.btn-group button.btn.btn-primary`).Each(func(i int, s *goquery.Selection) {
		direction := s.Find(`svg`).AttrOr("data-icon", "empty")
		if direction == "chevron-right" {
			paginateStatus = s.AttrOr("disabled", "not-disabled")
			fmt.Printf("%s", paginateStatus)
		}
	})

	return nil, map[string]interface{}{
		"pagination_status": paginateStatus,
		"episodes":          animeEpisodes,
	}
}

//get the episode info and there servers to play the anime video
func (k *kickass) EpisodesInfoSelector(serverCount int, path string) map[string]interface{} {
	server := k.getServerKey(serverCount)
	url := fmt.Sprintf("%s%s", server_engine.ParsedServers.KickassanimeServers[server], path)
	var innerHTML string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`footer#footer`),
		chromedp.InnerHTML(`#content div.player-wrapper.mb-3`, &innerHTML),
	}
	//run the assigned task
	InitializeChromeDp(func() chromedp.Tasks {
		return tasks
	})

	return map[string]interface{}{
		"innerHTML": innerHTML,
	}
}

//parse the given episode info scraped by the method epiusodeinfoselector
func (k *kickass) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	innerHTML := value["innerHTML"]
	byteData := []byte(innerHTML.(string))
	readerData := bytes.NewReader(byteData)

	//process the html and parse it to the main json value
	doc, err := goquery.NewDocumentFromReader(readerData)
	if err != nil {
		return utility.NewInternalServerError(fmt.Sprintln(err)), nil
	}

	//get the iframe src
	iframes := make([]string, 0)
	doc.Find(`div.ka-player.mb-3.embed-responsive.embed-responsive-16by9 iframe.embed-responsive-item`).Each(func(i int, s *goquery.Selection) {
		iframeUrl := s.AttrOr("src", "empty")
		fmt.Printf("%d->%s\n", i, iframeUrl)
		iframes = append(iframes, iframeUrl)
	})

	// div.player-wrapper.mb-3  div.ka-player.mb-3.embed-responsive.embed-responsive-16by9
	doc.Find(`div.player-wrapper.mb-3  div.ka-player.mb-3.embed-responsive.embed-responsive-16by9 select#ext-servers-select option`).Each(func(i int, s *goquery.Selection) {
		iframeUrl := s.AttrOr("value", "empty")
		fmt.Printf("%d->%s\n", i, iframeUrl)
		iframes = append(iframes, iframeUrl)
	})

	return nil, iframes
}

func (k *kickass) SearchAnimeSelector(keyword string) map[string]interface{} {
	return nil
}

func (k *kickass) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}
