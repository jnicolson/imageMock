package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

var (
	light_gray color.Color = color.RGBA{204, 204, 204, 255}
	dark_gray  color.Color = color.RGBA{150, 150, 150, 255}
)

type imageMock struct {
	font *truetype.Font
	dpi  float64
	size float64
}

func NewImageMock() *imageMock {
	return &imageMock{}
}

func (i *imageMock) setFont(fontfile string) {
	var fontBytes []byte
	var err error

	if fontfile == "internal" {
		fontBytes, _ = getData()
	} else {
		fontBytes, err = ioutil.ReadFile(fontfile)
		if err != nil {
			log.Fatal("Could not read font file")
		}
	}

	i.font, err = truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal("Could Not Parse Font Data")
	}
}

func (i *imageMock) setDpi(dpi float64) {
	i.dpi = dpi
}

func (i *imageMock) generateImage(x, y int) (*bytes.Buffer, error) {

	image_text := fmt.Sprintf("%d x %d", x, y)
	i.size = math.Max(
		math.Min(float64(x)/float64(len(image_text))*float64(0.7), float64(y)*float64(0.3)), 5)

	m := image.NewRGBA(image.Rect(0, 0, x, y))
	draw.Draw(m, m.Bounds(), &image.Uniform{light_gray}, image.ZP, draw.Over)

	c := freetype.NewContext()
	c.SetDPI(i.dpi)
	c.SetFont(i.font)
	c.SetClip(m.Bounds())
	c.SetDst(m)
	c.SetFontSize(i.size)
	c.SetSrc(&image.Uniform{dark_gray})

	x_pos := (m.Bounds().Max.X / 2) - int(i.getWidth(image_text)/2)
	y_pos := (m.Bounds().Max.Y / 2) + (i.getHeight() / 2)

	c.DrawString(image_text, freetype.Pt(x_pos, y_pos))

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, m, &jpeg.Options{Quality: 100}); err != nil {
		log.Println("Unable to encode image")
	}
	return buffer, nil
}

func (i *imageMock) getHeight() int {
	scale := fixed.Int26_6(i.size * i.dpi * (64.0 / 72.0))
	bounds := i.font.Bounds(scale)
	return int((bounds.Max.Y - bounds.Min.Y) >> 8)
}

func (i *imageMock) getWidth(s string) int {
	scale := fixed.Int26_6(i.size * i.dpi * (64.0 / 72.0))
	var width fixed.Int26_6

	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := i.font.Index(rune)
		if hasPrev {
			width += i.font.Kern(scale, prev, index) << 2
		}
		width += i.font.HMetric(scale, index).AdvanceWidth << 2
		prev, hasPrev = index, true
	}

	return int(width >> 8)
}
