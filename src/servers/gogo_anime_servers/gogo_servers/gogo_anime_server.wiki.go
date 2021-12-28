package gogo_servers

/*
This source file contains the scraping of gogoanimeserver.wiki
*/

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/egnimos/anime-scrapper/src/repository"
	"github.com/egnimos/anime-scrapper/src/server_engine"

	// "github.com/egnimos/anime-scrapper/src/servers"
	"github.com/egnimos/anime-scrapper/src/utility"
)

var (
	GetGogoAnimeServerWiki server_engine.AnimeServerInterface = &gogoAnimeWiki{}
)

type gogoAnimeWiki struct{}

// func (g *gogoAnimeWiki) getServerKey(serverCount int) string {
// 	serverName := "gogo_anime_server"
// 	return fmt.Sprintf("%s_%d", serverName, serverCount)
// }

// func (g *gogoAnimeWiki) getAnimeListingUrl(serverKey string, pageCount int) string {
// 	url := server_engine.ParsedServers.GogoanimeServers[serverKey]
// 	switch serverKey {
// 	case "gogo_anime_server_1": //https://gogoanime.cm
// 		return fmt.Sprintf("%s/anime-list.html?page=%d", url, pageCount)
// 	case "gogo_anime_server_2": //https://gogo-anime.su
// 		return fmt.Sprintf("%s/all_anime/?page=%d&alphabet=all", url, pageCount)
// 	case "gogo_anime_server_3": //https://ww2.gogoanimes.org
// 		return fmt.Sprintf("%s/anime-list?page=%d", url, pageCount)
// 	case "gogo_anime_server_4": //https://ww1.gogoanime2.org
// 		return fmt.Sprintf("%s/animelist/all/%d", url, pageCount)
// 	case "gogo_anime_server_5": //https://www1.gogoanime.sx
// 		return fmt.Sprintf("%s/anime-list?page=%d", url, pageCount)
// 	case "gogo_anime_server_6": //https://gogoanime.mom
// 		return fmt.Sprintf("%s/anime-list", url)
// 	case "gogo_anime_server_7": //https://gogoanime.wiki
// 		return fmt.Sprintf("%s/anime-list.html?page=%d", url, pageCount)
// 	default:
// 		return ""
// 	}
// }

func (g *gogoAnimeWiki) AnimeListingSelector(serverCount int, pageCount int) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerSx.AnimeListingSelector(serverCount, pageCount)
}

func (g *gogoAnimeWiki) AnimeListingHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeListings) {
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
		displayImageHTML := s.AttrOr("title", "")
		// splitedValue := strings.Split(displayImageHTML, "<")
		// posterImage := strings.Split(splitedValue[2], "\"")[1]
		// fmt.Println(posterImage)

		byteData := []byte(displayImageHTML)
		readerData := bytes.NewReader(byteData)

		//process the html and parse it to the main json value
		doc1, err := goquery.NewDocumentFromReader(readerData)
		if err != nil {
			log.Fatalln(err.Error())
		}

		posterImage := doc1.Find(`div.thumnail_tool img`).AttrOr("src", "")

		title := strings.TrimSpace(s.Text())
		navigationUrl := s.Find(`a`).AttrOr("href", "empty")
		fmt.Printf("%s\n%s\n", title, navigationUrl)
		animeListing := repository.AnimeListing{
			NavigationUrl:     navigationUrl,
			AnimeDisplayImage: posterImage,
			AnimeTitle:        title,
		}

		animeListings = append(animeListings, animeListing)
	})

	return nil, animeListings
}

//get the anime info
func (g *gogoAnimeWiki) AnimeInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.AnimeInfoSelector(serverCount, path)
}

//anime info html parser
func (g *gogoAnimeWiki) AnimeInfoHtmlParser(value map[string]interface{}) (utility.RestError, *repository.AnimeInfo) {
	return GetGogoAnimeServerCm.AnimeInfoHtmlParser(value)
}

//get the lsit of episodes
func (g *gogoAnimeWiki) EpisodesSelector(serverCount int, pageCount int, path string) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.EpisodesSelector(serverCount, pageCount, path)
}

//parse the episodes list
func (g *gogoAnimeWiki) EpisodesHtmlParser(value map[string]interface{}) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.EpisodesHtmlParser(value)
}

//get the episodes info
func (g *gogoAnimeWiki) EpisodesInfoSelector(serverCount int, path string) (utility.RestError, map[string]interface{}) {
	return GetGogoAnimeServerCm.EpisodesInfoSelector(serverCount, path)
}

//parse the episodes info
func (g *gogoAnimeWiki) EpisodesInfoHtmlParser(value map[string]interface{}) (utility.RestError, []string) {
	return GetGogoAnimeServerCm.EpisodesInfoHtmlParser(value)
}

func (g *gogoAnimeWiki) SearchAnimeSelector(keyword string) (utility.RestError, map[string]interface{}) {
	return nil, nil
}

func (g *gogoAnimeWiki) SearchAnimeHtmlParser(value map[string]interface{}) (utility.RestError, repository.AnimeInfos) {
	return nil, nil
}
