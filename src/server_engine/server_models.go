package server_engine

type ServerList struct {
	NineanimeServers    map[string]string `json:"9_anime_servers"`
	GogoanimeServers    map[string]string `json:"gogo_anime_servers"`
	KissanimeServers    map[string]string `json:"kiss_anime_servers"`
	AnimepaheServers    map[string]string `json:"animepahe_servers"`
	AnimekisaServers    map[string]string `json:"animekisa_servers"`
	TwistmoeServers     map[string]string `json:"twist_moe_servers"`
	AnimeheavenServers  map[string]string `json:"animeheaven_servers"`
	KickassanimeServers map[string]string `json:"kickass_anime_servers"`
	FouranimeServers    map[string]string `json:"4_anime_servers"`
}

