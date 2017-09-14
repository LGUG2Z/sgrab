package sonarr_test

import (
	"fmt"
	"net/http"

	. "github.com/lgug2z/sgrab/sonarr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Sonarr", func() {
	var server *ghttp.Server
	var sonarr Client

	series := []Series{{Title: "Westworld", ID: 1}}
	episodes := []Episode{{Title: "The Original", ID: 1, EpisodeFileID: 1, SeriesID: 1, SeasonNumber: 1, EpisodeNumber: 1}}
	episodeFile := EpisodeFile{ID: 1, Path: "/path/to/episode/1.mkv"}

	BeforeEach(func() {
		server = ghttp.NewServer()
		sonarr.URL = server.URL()
		sonarr.Client = http.Client{}
	})

	Describe("When looking up series on a valid server", func() {
		It("Returns a list of Series objects", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/series/")),
					ghttp.RespondWithJSONEncoded(http.StatusOK, series),
				),
			)

			series, err := sonarr.Series()
			Expect(err).ToNot(HaveOccurred())
			Expect(series).To(Equal([]Series{{Title: "Westworld", ID: 1}}))
		})
	})

	Describe("When looking up episodes of a valid series on the server", func() {
		It("Returns a list of Episode objects for that series", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/episode/")),
					ghttp.VerifyFormKV("seriesId", "1"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, episodes),
				),
			)
			episodes, err := sonarr.Episodes(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(episodes).To(Equal([]Episode{
				{Title: "The Original", ID: 1, EpisodeFileID: 1, SeriesID: 1, SeasonNumber: 1, EpisodeNumber: 1},
			}))
		})
	})

	Describe("When looking up episodes of a valid series on the server", func() {
		It("Returns a list of Episode objects for that series", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/episodeFile/1")),
					ghttp.RespondWithJSONEncoded(http.StatusOK, episodeFile),
				),
			)
			episodes, err := sonarr.EpisodeFile(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(episodes).To(Equal(EpisodeFile{ID: 1, Path: "/path/to/episode/1.mkv"}))
		})
	})

	Describe("When making a request with an invalid API key", func() {
		It("An error is returned", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/series/")),
					ghttp.RespondWith(http.StatusUnauthorized, ""),
				),
			)
			_, err := sonarr.Series()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(ErrUnauthorized.Error()))
		})
	})
})
