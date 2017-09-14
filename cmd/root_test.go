package cmd_test

import (
	. "github.com/lgug2z/sgrab/cmd"

	"net/http"

	"fmt"

	"os"

	"github.com/lgug2z/sgrab/sonarr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/spf13/afero"
)

var _ = Describe("SGrab", func() {
	Describe("When run without the required flags", func() {
		It("Should return an error", func() {
			err := SGrab(nil, Flags{}, sonarr.Client{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(ErrInformationMissing.Error()))

		})
	})

	Describe("When a requested series does not exist on the seedbox", func() {
		It("Should return an error", func() {
			server := ghttp.NewServer()
			sonarr := sonarr.Client{}
			sonarr.URL = server.URL()
			sonarr.Client = http.Client{}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/series/")),
					ghttp.RespondWith(http.StatusOK, "[]"),
				),
			)

			f := Flags{
				APIKey:     "aaa",
				Episode:    "s01e01",
				SeedboxURL: "ddd",
				Series:     "Not a Real Series",
				SonarrURL:  "bbb",
				Username:   "ccc",
			}

			err := SGrab(nil, f, sonarr)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(ErrCouldNotFindSeries("Not a Real Series").Error()))

		})
	})

	Describe("When called for a valid file", func() {
		series := []sonarr.Series{{Title: "Westworld", ID: 1}}
		episodes := []sonarr.Episode{{Title: "The Original", ID: 1, EpisodeFileID: 1, SeriesID: 1, SeasonNumber: 1, EpisodeNumber: 1}}
		episodeFile := sonarr.EpisodeFile{ID: 1, Path: "/westworld-s01e01.mkv"}
		f := Flags{
			APIKey:         "key",
			Episode:        "s01e01",
			SSHKeyLocation: fmt.Sprintf("%s/.vagrant.d/insecure_private_key", os.Getenv("HOME")),
			SeedboxURL:     "127.0.0.1",
			Series:         "Westworld",
			Username:       "vagrant",
			Port:           "2222",
		}

		It("Should return an error if the SFTP credentials are not correct", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/series/")),
					ghttp.RespondWithJSONEncoded(http.StatusOK, series),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/episode/")),
					ghttp.RespondWithJSONEncoded(http.StatusOK, episodes),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/episodeFile/1")),
					ghttp.RespondWithJSONEncoded(http.StatusOK, episodeFile),
				),
			)

			f.SonarrURL = server.URL()
			f.Username = "wrong"

			sonarr := sonarr.Client{
				URL:    server.URL(),
				APIKey: "key",
				Client: http.Client{},
			}

			err = SGrab(afero.NewMemMapFs(), f, sonarr)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(ErrCredentialsRejected.Error()))
		})

		It("Should transfer the file via SFTP if credentials are correct", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/series/")),
					ghttp.RespondWithJSONEncoded(http.StatusOK, series),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/episode/")),
					ghttp.RespondWithJSONEncoded(http.StatusOK, episodes),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("/api/episodeFile/1")),
					ghttp.RespondWithJSONEncoded(http.StatusOK, episodeFile),
				),
			)

			f.SonarrURL = server.URL()
			f.Username = "vagrant"

			sonarr := sonarr.Client{
				URL:    server.URL(),
				APIKey: "key",
				Client: http.Client{},
			}

			Expect(SGrab(afero.NewMemMapFs(), f, sonarr)).To(Succeed())
		})
	})

})
