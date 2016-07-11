//go:generate cat static/index.html

package main

import (
	"fmt"
	m "jarl/middleware"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
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

	router := httprouter.New()
	router.GET("/", mainHandler)
	router.GET("/:size", imageHandler)

	loggedRouter := m.LogMiddleware(router)

	log.Printf("Starting to serve on port %s", *port)
	err := http.ListenAndServe(":"+*port, loggedRouter)

	if err != nil {
		log.Fatal("Listen and Serve error: ")
	}

}

func mainHandler(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	params := []httprouter.Param{}
	param := httprouter.Param{Key: "size", Value: "800x600"}
	params = append(params, param)

	imageHandler(rw, req, params)
}

func imageHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	size := params.ByName("size")

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
		fmt.Fprintf(rw, "Could not parse size")
		return
	}

	image, err := im.generateImage(x, y)

	if err != nil {
		log.Fatal("Could not generate image")
	}

	rw.Header().Set("Content-Type", "image/jpeg")
	rw.Header().Set("Content-Length", strconv.Itoa(len(image.Bytes())))
	if _, err := rw.Write(image.Bytes()); err != nil {
		log.Println("Unable to write image")
	}

}
