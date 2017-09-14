package sonarr

import "time"

type EpisodeFile struct {
	SeriesID     int       `json:"seriesId"`
	SeasonNumber int       `json:"seasonNumber"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	DateAdded    time.Time `json:"dateAdded"`
	SceneName    string    `json:"sceneName"`
	Quality      struct {
		Quality struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"quality"`
		Proper bool `json:"proper"`
	} `json:"quality"`
	ID int `json:"id"`
}
