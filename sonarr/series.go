package sonarr

import "time"

type Series struct {
	Title           string `json:"title"`
	AlternateTitles []struct {
		Title        string `json:"title"`
		SeasonNumber int    `json:"seasonNumber"`
	} `json:"alternateTitles"`
	SortTitle         string    `json:"sortTitle"`
	SeasonCount       int       `json:"seasonCount"`
	TotalEpisodeCount int       `json:"totalEpisodeCount"`
	EpisodeCount      int       `json:"episodeCount"`
	EpisodeFileCount  int       `json:"episodeFileCount"`
	SizeOnDisk        int64     `json:"sizeOnDisk"`
	Status            string    `json:"status"`
	Overview          string    `json:"overview"`
	PreviousAiring    time.Time `json:"previousAiring"`
	Network           string    `json:"network"`
	AirTime           string    `json:"airTime"`
	Images            []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	Seasons []struct {
		SeasonNumber int  `json:"seasonNumber"`
		Monitored    bool `json:"monitored"`
		Statistics   struct {
			PreviousAiring    time.Time `json:"previousAiring"`
			EpisodeFileCount  int       `json:"episodeFileCount"`
			EpisodeCount      int       `json:"episodeCount"`
			TotalEpisodeCount int       `json:"totalEpisodeCount"`
			SizeOnDisk        int64     `json:"sizeOnDisk"`
			PercentOfEpisodes float64   `json:"percentOfEpisodes"`
		} `json:"statistics"`
	} `json:"seasons"`
	Year              int           `json:"year"`
	Path              string        `json:"path"`
	ProfileID         int           `json:"profileId"`
	SeasonFolder      bool          `json:"seasonFolder"`
	Monitored         bool          `json:"monitored"`
	UseSceneNumbering bool          `json:"useSceneNumbering"`
	Runtime           int           `json:"runtime"`
	TvdbID            int           `json:"tvdbId"`
	TvRageID          int           `json:"tvRageId"`
	TvMazeID          int           `json:"tvMazeId"`
	FirstAired        time.Time     `json:"firstAired"`
	LastInfoSync      time.Time     `json:"lastInfoSync"`
	SeriesType        string        `json:"seriesType"`
	CleanTitle        string        `json:"cleanTitle"`
	ImdbID            string        `json:"imdbId"`
	TitleSlug         string        `json:"titleSlug"`
	Certification     string        `json:"certification"`
	Genres            []string      `json:"genres"`
	Tags              []interface{} `json:"tags"`
	Added             time.Time     `json:"added"`
	Ratings           struct {
		Votes int     `json:"votes"`
		Value float64 `json:"value"`
	} `json:"ratings"`
	QualityProfileID int `json:"qualityProfileId"`
	ID               int `json:"id"`
}
