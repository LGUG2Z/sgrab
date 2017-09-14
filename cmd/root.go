package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"net/http"

	"crypto/tls"

	"path/filepath"

	"github.com/lgug2z/sgrab/sonarr"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

var RootCmd = &cobra.Command{
	Use:   "sgrab",
	Short: "sgrab is a command line utility to grab episodes from a seedbox running Sonarr.",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		sshKey, err := GetKeyFile(rootFlags.SSHKeyLocation)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Allow self signed certs
		sonarrClient := sonarr.Client{
			APIKey: rootFlags.APIKey,
			URL:    urlWithSlash(rootFlags.SonarrURL),
			Client: http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			},
		}

		fs := afero.NewOsFs()

		if err := SGrab(fs, rootFlags, sshKey, sonarrClient); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func SGrab(fs afero.Fs, f Flags, k ssh.Signer, c sonarr.SonarrClient) error {
	if !hasRequiredFlags(f) {
		return ErrInformationMissing
	}

	series, err := c.Series()
	if err != nil {
		return err
	}

	requestedSeries, err := findSeries(series, f.Series)
	if err != nil {
		return err
	}

	episodes, err := c.Episodes(requestedSeries.ID)
	if err != nil {
		return err
	}

	episodeFileID, err := findEpisodeFileID(episodes, f.Episode)
	if err != nil {
		return err
	}

	episodeFile, err := c.EpisodeFile(episodeFileID)
	if err != nil {
		return err
	}

	// Set the destination path to the present working directory
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dstPath := fmt.Sprintf("%s/%s", pwd, filepath.Base(episodeFile.Path))

	// Starting copying the file in a new goroutine
	copyChan := make(chan error)
	go func() {
		copyChan <- copyFile(fs, f, k, episodeFile.Path, dstPath)
	}()

	// Listen for an interrupt signal in a new goroutine
	signalChan := make(chan os.Signal, 1)
	cleanupChan := make(chan error)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			// Cleanup an incomplete file transfer
			cleanupChan <- fs.Remove(dstPath)
		}
	}()

	// Block and wait for either
	select {
	case err := <-copyChan:
		// 1) The copy function to return either nil or an error
		return err
	case err = <-cleanupChan:
		// 2) An interrupt signal to be received and trigger the cleanup
		if err == nil {
			// And for the cleanup to be either successful
			return ErrInterruptReceived
		}
		// Or unsuccessful
		return ErrInterruptReceivedCleanupFailed
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootFlags Flags

func init() {
	viper.SetEnvPrefix("sgrab")
	viper.AutomaticEnv()

	RootCmd.Flags().StringVar(&rootFlags.SonarrURL, "sonarr", viper.GetString("sonarr_url"), "Sonarr url")
	RootCmd.Flags().StringVar(&rootFlags.APIKey, "api-key", viper.GetString("api_key"), "Sonarr API key")
	RootCmd.Flags().StringVarP(&rootFlags.Series, "series", "s", "", "Series name")
	RootCmd.Flags().StringVarP(&rootFlags.Episode, "episode", "e", "", "Episode number (format \"s01e02\")")
	RootCmd.Flags().StringVar(&rootFlags.SeedboxURL, "seedbox", viper.GetString("seedbox"), "Seedbox address")
	RootCmd.Flags().StringVar(&rootFlags.Username, "username", viper.GetString("username"), "Seedbox login username")
	RootCmd.Flags().StringVar(&rootFlags.SSHKeyLocation, "ssh-key", fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME")), "Path to SSH key")
	RootCmd.Flags().StringVar(&rootFlags.Port, "port", "22", "SSH port number for seedbox")
}
