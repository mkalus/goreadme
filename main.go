package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/Depado/bfchroma"
	bf "github.com/russross/blackfriday/v2"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// embed css file
//go:embed md.css
var css []byte

// variables set by flags
var index string
var sourcePath string
var serverAddress string
var help bool
var logging bool
var quiet bool

// set by goreleaser
var version = "snapshot"
var commit = "unknown"
var date = "unknown"

// initialize flags
func init() {
	flag.StringVar(&index, "index", "README.md", "index file")
	flag.StringVar(&sourcePath, "source", ".", "source path (relative or absolute")
	flag.StringVar(&serverAddress, "address", ":8080", "server address to listen to")
	flag.BoolVar(&logging, "logging", true, "turn logging on or off")
	flag.BoolVar(&quiet, "quiet", false, "do not print header when starting program")
	flag.BoolVar(&help, "h", false, "print help (shorthand)")
	flag.BoolVar(&help, "h", false, "print help (shorthand)")
	flag.Parse()
}

func printHeader() {
	// shorten commit
	if len(commit) > 7 {
		commit = commit[0:8]
	}

	println("goreadme - Very simple README.md viewer written in go")
	println("Version: " + version + ", commit: " + commit + ", built at " + date)
}

func main() {
	if !quiet {
		printHeader()
	}

	if logging {
		log.Printf("Starting server at address: %s", serverAddress)
	}

	// just print help and exit
	if help {
		println()
		flag.PrintDefaults()
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// remove starting slash
		if path[0] == '/' {
			path = path[1:]
		}
		// empty path -> load index
		if path == "" {
			path = index
		}

		// serve css file? serve embedded code
		if path == "md.css" {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			_, _ = w.Write(css)
			return
		}

		// add source path
		path = filepath.Join(sourcePath, path)

		// find file
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			if logging {
				log.Printf("404: %s does not exist", path)
			}
			handleHTTPError(w, 404)
			return
		}
		if err != nil {
			if logging {
				log.Printf("500: file error occurred: %s", err)
			}
			handleHTTPError(w, 500)
			return
		}
		if info.IsDir() {
			if logging {
				log.Printf("500: %s is a path", path)
			}
			handleHTTPError(w, 500)
			return
		}

		// get extension of file
		extension := filepath.Ext(path)
		if extension == ".md" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				if logging {
					log.Printf("500: read error occurred: %s", err)
				}
				handleHTTPError(w, 500)
				return
			}

			// write header
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <title>`+html.EscapeString(path)+`</title>
  <link rel="stylesheet" href="/md.css">
</head>
<body><div class="container">
`)

			result := render(data)
			_, _ = w.Write(result)

			_, _ = fmt.Fprint(w, `</div></body></html>
`)
			if logging {
				log.Printf("200: markdown: %s", path)
			}
			return
		}

		// serve other files
		http.ServeFile(w, r, path)
		if logging {
			log.Printf("200: served: %s", path)
		}
	})

	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

// Defines the extensions that are used
var exts = bf.NoIntraEmphasis | bf.Tables | bf.FencedCode | bf.Autolink |
	bf.Strikethrough | bf.SpaceHeadings | bf.BackslashLineBreak |
	bf.DefinitionLists | bf.Footnotes

// Defines the HTML rendering flags that are used
var flags = bf.Smartypants | bf.SmartypantsFractions | bf.SmartypantsDashes | bf.SmartypantsLatexDashes | bf.TOC |
	bf.FootnoteReturnLinks | bf.HrefTargetBlank | bf.CompletePage

// handle error by writing header and info about error to response
func handleHTTPError(w http.ResponseWriter, error int) {
	w.WriteHeader(error)
	var statusLine string
	switch error {
	case 404:
		statusLine = "404 not found"
	default:
		statusLine = "500 internal server error"
	}
	_, _ = fmt.Fprint(w, statusLine)
}

// render will take a []byte input and will render it using a new renderer each
// time because reusing the same can mess with TOC and header IDs
func render(input []byte) []byte {
	return bf.Run(
		input,
		bf.WithRenderer(
			bfchroma.NewRenderer(
				bfchroma.Extend(
					bf.NewHTMLRenderer(bf.HTMLRendererParameters{
						Flags: flags,
					}),
				),
			),
		),
		bf.WithExtensions(exts),
	)
}
