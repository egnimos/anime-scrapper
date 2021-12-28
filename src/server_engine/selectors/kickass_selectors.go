package selectors

var (
	GetKickAssSelectors KickAssSelectors = &kickAssSelectors{}
)

type KickAssSelectors interface {
	AnimeListing
	AnimeInfo
	Episodes
	EpisodesInfo
}

type kickAssSelectors struct{}

/*
This section contains the selectors for fetching and processing
the Anime Listing
*/
type AnimeListing interface {
	AnimeListingWaitVisible() string
	AnimeListingSwitchToGrid() string
	AnimeListingPaginate() string
	AnimeListingInnerHTML() string
	AnimeInfoForAnimeListing() string
	AnimeListingNavigationUrl() string
	AnimeListingBgUrl() string
	AnimeListingTitle() string
}

func (ka *kickAssSelectors) AnimeListingWaitVisible() string {
	return `footer#footer`
}

func (ka *kickAssSelectors) AnimeListingSwitchToGrid() string {
	return `input#__BVID__50__BV_option_1_`
}

func (ka *kickAssSelectors) AnimeListingPaginate() string {
	return `#content button.btn.btn-primary svg.svg-inline--fa.fa-chevron-right.fa-w-10`
}

func (ka *kickAssSelectors) AnimeListingInnerHTML() string {
	return `#content div.row.video-list.row.mx-0 div.col.row.video-list.row.mx-0`
}

func (ka *kickAssSelectors) AnimeInfoForAnimeListing() string {
	return `div.video-item.col-6.mb-2.px-1.col-2`
}

func (ka *kickAssSelectors) AnimeListingNavigationUrl() string {
	return `a.ka-url-wrapper.video-item-poster.rounded`
}

func (ka *kickAssSelectors) AnimeListingBgUrl() string {
	return ka.AnimeListingNavigationUrl()
}

func (ka *kickAssSelectors) AnimeListingTitle() string {
	return `a.ka-url-wrapper.video-item-title`
}

/**
Anime Info this section contains the selectors
for kickass anime anime info
**/
type AnimeInfo interface {
	AnimeInfoWaitVisible() string
	AnimeInfoInnerHTML() string
	AnimeInfo1stBlock() string
	AnimeInfo1stBlockAnimeImage() string
	AnimeInfo1stBlockAnimeMap() string
	Anime1stBlockKey() string
	Anime1stBlockValue() string
	AnimeInfo2ndBlock() string
	AnimeInfo2ndBlockTitle() string
	AnimeInfo2ndBlockSummary() string
	AnimeInfo2ndBlockGenres() string
}

func (ka *kickAssSelectors) AnimeInfoWaitVisible() string {
	return ka.AnimeListingWaitVisible()
}

func (ka *kickAssSelectors) AnimeInfoInnerHTML() string {
	return `div#main.container-fluid.pt-3`
}

func (ka *kickAssSelectors) AnimeInfo1stBlock() string {
	return `div#sidebar-anime-info.border.rounded.mb-3`
}

func (ka *kickAssSelectors) AnimeInfo1stBlockAnimeImage() string {
	return `div.poster`
}

func (ka *kickAssSelectors) AnimeInfo1stBlockAnimeMap() string {
	return `div.p-3 div.mb-2`
}

func (ka *kickAssSelectors) Anime1stBlockKey() string {
	return `div.font-weight-bold.mr-1`
}

func (ka *kickAssSelectors) Anime1stBlockValue() string {
	return `div.mb-2 span`
}

func (ka *kickAssSelectors) AnimeInfo2ndBlock() string {
	return `#content div.anime-info.border.rounded.mb-3`
}

func (ka *kickAssSelectors) AnimeInfo2ndBlockTitle() string {
	return `a.ka-url-wrapper div.info-header.p-3.px-4.hep div.info-wrapper h1.title`
}

func (ka *kickAssSelectors) AnimeInfo2ndBlockSummary() string {
	return `div.container.p-3 div.mb-3 p.mb-0`
}

func (ka *kickAssSelectors) AnimeInfo2ndBlockGenres() string {
	return `div.mb-3 a.ka-url-wrapper.d-inline-block`
}

/*
This section contains the selectors for Fetching and Processing
the list of Episodes
*/
type Episodes interface {
	EpisodesWaitVisible() string
	EpisodesPaginate() string
	EpisodesRearrange() string
	EpisodesInnerHTML() string
	EpisodesInfoFind() string
	EpisodesPathUrl() string
	EpisodeName() string
	EpisodesPaginationStatus() string
}

func (ka *kickAssSelectors) EpisodesWaitVisible() string {
	return ka.AnimeInfoWaitVisible()
}

func (ka *kickAssSelectors) EpisodesPaginate() string {
	return `#content div.main-episode-list.border.rounded.p-3.mb-3 button.btn.btn-primary`
}

func (ka *kickAssSelectors) EpisodesRearrange() string {
	return `#content div.main-episode-list.border.rounded.p-3.mb-3 table.table.b-table.table-hover thead[role="rowgroup"] tr[role="row"] th[aria-colindex="1"] span.sr-only`
}

func (ka *kickAssSelectors) EpisodesInnerHTML() string {
	return `#content div.main-episode-list.border.rounded.p-3.mb-3`
}

func (ka *kickAssSelectors) EpisodesInfoFind() string {
	return `table.table.b-table.table-hover tbody[role="rowgroup"] tr[role="row"]`
}

func (ka *kickAssSelectors) EpisodesPathUrl() string {
	return `td[aria-colindex="1"] a.ka-url-wrapper`
}

func (ka *kickAssSelectors) EpisodeName() string {
	return ka.EpisodesPathUrl()
}

func (ka *kickAssSelectors) EpisodesPaginationStatus() string {
	return `div.row.align-items-end div.text-right.col-lg-3.col-12 div.mb-3.btn-group button.btn.btn-primary`
}

/*
This section contains the selectors for Fetching and Processing
the Episode info
*/
type EpisodesInfo interface {
	EpisodesInfoWaitVisible() string
	EpisodesInfoInnerHTML() string
	EpisodesInfo1stBlockVideoServers() string
	EpisodesInfo2ndBlockVideoServers() string
}

func (ka *kickAssSelectors) EpisodesInfoWaitVisible() string {
	return ka.EpisodesWaitVisible()
}

func (ka *kickAssSelectors) EpisodesInfoInnerHTML() string {
	return `#content div.player-wrapper.mb-3`
}

func (ka *kickAssSelectors) EpisodesInfo1stBlockVideoServers() string {
	return `div.ka-player.mb-3.embed-responsive.embed-responsive-16by9 iframe.embed-responsive-item`
}

func (ka *kickAssSelectors) EpisodesInfo2ndBlockVideoServers() string {
	return `div.player-wrapper.mb-3  div.ka-player.mb-3.embed-responsive.embed-responsive-16by9 select#ext-servers-select option`
}
