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
