# GoReadme

Very simple README.md viewer written in go.

Just start binary in a folder containing a README.md file and be set.

Other features:

* Serves static files from folder/subfolders.
* Automatic code highlighting
* Many md extensions activated

## Parameters

* `-help` or `-h`: print help
* `-index`: Define index file (default: `README.md`)
* `-source`: Source path, may be realtive or absolute (default: `.`)
* `-address`: Address to listen to (default `:8080`)
* `-logging`: Turn logging on or off (default `true`)
* `-quiet`: Do not print header when starting goreadme

## Docker

Docker image is available:

```bash
docker pull ronix/goreadme
```

Images are available for amd64, arm64 and 386.

## Snap

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-white.svg)](https://snapcraft.io/goreadme)

Snaps are available, too. Install using:

```bash
sudo snap install goreadme
```

**Note:** Due to snap security restrictions, goreadme snap can only access files in the user's home folder, so make sure
the binary runs with proper permissions. Otherwise you will get http 500 errors and the log will show access denied
entries.
