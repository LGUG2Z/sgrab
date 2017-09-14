package sonarr

import "time"

type Episode struct {
	SeriesID              int       `json:"seriesId"`
	EpisodeFileID         int       `json:"episodeFileId"`
	SeasonNumber          int       `json:"seasonNumber"`
	EpisodeNumber         int       `json:"episodeNumber"`
	Title                 string    `json:"title"`
	AirDate               string    `json:"airDate"`
	AirDateUtc            time.Time `json:"airDateUtc"`
	Overview              string    `json:"overview"`
	HasFile               bool      `json:"hasFile"`
	Monitored             bool      `json:"monitored"`
	SceneEpisodeNumber    int       `json:"sceneEpisodeNumber"`
	SceneSeasonNumber     int       `json:"sceneSeasonNumber"`
	TvDbEpisodeID         int       `json:"tvDbEpisodeId"`
	AbsoluteEpisodeNumber int       `json:"absoluteEpisodeNumber"`
	ID                    int       `json:"id"`
}
