package sonarr

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"errors"

	"github.com/dghubble/sling"
)

const APIEndpoint = "api/"
const SeriesEndpoint = "series/"
const EpisodeFileEndpoint = "episodeFile/"
const EpisodeEndpoint = "episode/"

var ErrUnauthorized = errors.New("API key was rejected by Sonarr.")

type SonarrClient interface {
	Series() ([]Series, error)
	Episodes(seriesID int) ([]Episode, error)
	EpisodeFile(episodeFileID int) (EpisodeFile, error)
}

type Client struct {
	APIKey string
	Client http.Client
	URL    string
}

func (c Client) Series() ([]Series, error) {
	var series []Series

	req, err := sling.
		New().
		Get(c.URL).
		Path(APIEndpoint).
		Path(SeriesEndpoint).
		Set("X-Api-Key", c.APIKey).
		Request()

	if err != nil {
		return series, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return series, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return series, ErrUnauthorized
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return series, err
	}

	if err := json.Unmarshal(bytes, &series); err != nil {
		return series, err
	}

	return series, nil
}

func (c Client) Episodes(seriesID int) ([]Episode, error) {
	var episodes []Episode

	req, err := sling.
		New().
		Get(c.URL).
		Path(APIEndpoint).
		Path(EpisodeEndpoint).
		QueryStruct(struct {
			SeriesID int `url:"seriesId,omitempty"`
		}{SeriesID: seriesID}).
		Set("X-Api-Key", c.APIKey).
		Request()

	if err != nil {
		return episodes, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return episodes, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return episodes, ErrUnauthorized
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return episodes, err
	}

	if err := json.Unmarshal(bytes, &episodes); err != nil {
		return episodes, err
	}

	return episodes, nil
}

func (c Client) EpisodeFile(episodeFileID int) (EpisodeFile, error) {
	var episodeFile EpisodeFile

	req, err := sling.
		New().
		Get(c.URL).
		Path(APIEndpoint).
		Path(EpisodeFileEndpoint).
		Path(strconv.Itoa(episodeFileID)).
		Set("X-Api-Key", c.APIKey).
		Request()

	if err != nil {
		return episodeFile, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return episodeFile, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return episodeFile, ErrUnauthorized
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return episodeFile, err
	}

	if err := json.Unmarshal(bytes, &episodeFile); err != nil {
		return episodeFile, err
	}

	return episodeFile, nil
}
