//go:generate cat static/index.html

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/namsral/flag"
)

// flags
var (
	dpi      = flag.Float64("dpi", 96, "Screen resolution in dots per inch")
	fontfile = flag.String("fontfile", "internal", "Filename of the ttf font file")
	port     = flag.String("port", "8080", "HTTP Port to listen on")
)

var im *imageMock

func main() {
	flag.Parse()

	im = NewImageMock()
	im.setFont(*fontfile)
	im.setDpi(*dpi)

	router := chi.NewRouter()
	router.Get("/", Index)
	router.Get("/{size}", imageHandler)

	loggedRouter := LogMiddleware(router)

	log.Printf("Starting to serve on port %s", *port)
	err := http.ListenAndServe(":"+*port, loggedRouter)

	if err != nil {
		log.Fatal("Listen and Serve error: ")
	}

}

func Index(w http.ResponseWriter, r *http.Request) {
	rctx := chi.RouteContext(r.Context())
	rctx.URLParams.Add("size", "800x600")

	imageHandler(w, r)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	size := chi.URLParam(r, "size")

	if size == "favicon.ico" {
		return
	}

	dimensions := strings.Split(size, "x")
	var x, y int
	var err error
	if len(dimensions) > 1 {
		x, err = strconv.Atoi(dimensions[0])
		y, err = strconv.Atoi(dimensions[1])
	} else {
		x, err = strconv.Atoi(dimensions[0])
		y, err = strconv.Atoi(dimensions[0])
	}

	if err != nil {
		fmt.Fprintf(w, "Could not parse size")
		return
	}

	image, err := im.generateImage(x, y)

	if err != nil {
		log.Fatal("Could not generate image")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(image.Bytes())))
	if _, err := w.Write(image.Bytes()); err != nil {
		log.Println("Unable to write image")
	}

}
