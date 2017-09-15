# sgrab

sgrab is a command line utility for grabbing episodes of TV shows from a 
seedbox running [Sonarr](https://sonarr.tv/).

## Requirements
* [Go](https://github.com/golang/go)

## Install
The latest version of `sgrab` can be installed using `go get`.

```
go get -u github.com/LGUG2Z/sgrab
```

Make sure that `$GOPATH` is set correctly that and that `$GOPATH/bin` is in your `$PATH`.

The `sgrab` executable will be installed under the `$GOPATH/bin` directory.

### Required information
In order to use `sgrab` the following information is required:

* Sonarr URL
* Sonarr API key
* Seedbox address
* Seedbox login username

This information can either be set in your shell rc:

```bash
export SGRAB_SONARR=https://mybox.com/sonarr/
export SGRAB_API_KEY=xxx
export SGRAB_SEEDBOX=mybox.com
export SGRAB_USERNAME=xxx
```

Or set using the relevant flags
```bash
Flags:
      --api-key string    Sonarr API key
  -e, --episode string    Episode number
  -h, --help              help for sgrab
      --port string       SSH port number for seedbox
      --seedbox string    Seedbox address
  -s, --series string     Series name
      --sonarr string     Sonarr url
      --ssh-key string    Path to SSH key
      --username string   Seedbox login username
```

## Usage
The key at $HOME/.ssh/id_rsa is used to establish a secure connection to the
seedbox to download the file. A different key can be provided using the `--ssh-key`
flag.

sgrab will by default try to connect to the seedbox on port 22. An alternative
port can be specified using the `--port` flag.

The `--series` flag is case-insensitive, however the name of the series must
otherwise match the primary name given to a series by Sonarr. Series that have
multi-word titles should be quoted.

The `--episode` flag uses the format "s01e02", with mandatory leading zeroes.

Example:

```bash
sgrab --series "Terrace House: Boys x Girls Next Door" --episode s01e01
```

