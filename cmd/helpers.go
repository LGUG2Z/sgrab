package cmd

import (
	"io"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/lgug2z/sgrab/sonarr"
	"github.com/pkg/sftp"
	pb "gopkg.in/cheggaaa/pb.v1"

	"os/user"

	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"golang.org/x/crypto/ssh"
)

func hasRequiredFlags(f Flags) bool {
	return len(f.SeedboxURL) > 0 && len(f.SonarrURL) > 0 && len(f.APIKey) > 0 && len(f.Username) > 0
}

type Flags struct {
	APIKey         string
	Episode        string
	SSHKeyLocation string
	SeedboxURL     string
	Series         string
	SonarrURL      string
	Username       string
	Port           string
}

func urlWithSlash(url string) string {
	if !strings.HasSuffix(url, "/") {
		url = fmt.Sprintf("%s/", url)
	}

	return url
}

func GetKeyFile(location string) (key ssh.Signer, err error) {
	if len(location) < 1 {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		location = usr.HomeDir + "/.ssh/id_rsa"
	}

	buf, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}

	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// From https://stackoverflow.com/questions/45441735/ssh-handshake-complains-about-missing-host-key
func getHostKey(fs afero.Fs, host string) (ssh.PublicKey, error) {
	file, err := fs.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return nil, fmt.Errorf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		return nil, fmt.Errorf("no hostkey for %s", host)
	}

	return hostKey, nil
}

func copyFile(fs afero.Fs, f Flags, k ssh.Signer, e, dstPath string) error {
	// Make sure there is an entry for the seedbox in $HOME/.ssh/known_hosts before connecting
	// Comment out when running test with Vagrant
	hostKey, err := getHostKey(fs, f.SeedboxURL)
	if err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: f.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(k),
		},
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use this when running tests with Vagrant
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// Make an SSH connection
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", f.SeedboxURL, f.Port), sshConfig)
	if err != nil {
		return ErrCredentialsRejected
	}
	defer client.Close()

	// Open an SFTP session over the SSH connection
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftp.Close()

	// Open the episode file on the seedbox
	src, err := sftp.Open(e)
	if err != nil {
		return err
	}
	defer src.Close()

	// Get the episode file info
	fi, err := src.Stat()
	if err != nil {
		return err
	}

	// Create the local file to copy to
	dst, err := fs.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Start a status bar based on the episode file size
	fmt.Println(fi.Name())
	bar := pb.New64(fi.Size()).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	bar.Start()
	defer bar.Finish()

	// Create a MultiWriter to write to the local file and send info to the progress bar
	dstWriter := io.MultiWriter(dst, bar)

	// Copy the file
	_, err = io.Copy(dstWriter, src)
	if err != nil {
		return err
	}

	return nil
}

func findEpisodeFileID(episodes []sonarr.Episode, toFind string) (int, error) {
	toFind = strings.ToLower(toFind)

	seasonRegex := regexp.MustCompile(`s\d{2}`)
	s0x := seasonRegex.FindString(toFind)

	episodeRegex := regexp.MustCompile(`e\d{2}`)
	e0x := episodeRegex.FindString(toFind)

	season, err := strconv.Atoi(strings.TrimPrefix(s0x, "s"))
	if err != nil {
		return 0, err
	}

	episode, err := strconv.Atoi(strings.TrimPrefix(e0x, "e"))
	if err != nil {
		return 0, err
	}

	for _, e := range episodes {
		if e.SeasonNumber == season && e.EpisodeNumber == episode {
			return e.EpisodeFileID, nil
		}
	}

	return 0, nil

}

func findSeries(series []sonarr.Series, toFind string) (sonarr.Series, error) {
	for _, s := range series {
		if strings.ToLower(s.Title) == strings.ToLower(toFind) {
			return s, nil
		}
	}

	return sonarr.Series{}, ErrCouldNotFindSeries(toFind)
}
