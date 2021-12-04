package repository

type AnimeInfo struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Poster  string `json:"poster"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Generes string `json:"generes"`
}

type AnimeListing struct {
	NavigationUrl     string `json:"navigation_url"`
	AnimeDisplayImage string `json:"anime_display_image"`
	AnimeTitle        string `json:"anime_title"`
}

type AnimeEpisode struct {
	EpisodePath string `json:"episode_path"`
	Episode     string `json:"episode"`
}

type AnimeListings []AnimeListing
type AnimeInfos []AnimeInfo
type AnimeEpisodes []AnimeEpisode
