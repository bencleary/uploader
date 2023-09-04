package uploader

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

type ImageFormat interface {
	Decode(file *os.File) (image.Image, error)
	Encode(file *os.File, img image.Image) error
}

type JPEGFormat struct{}
type PNGFormat struct{}
type GIFFormat struct{}

func (j *JPEGFormat) Decode(file *os.File) (image.Image, error) {
	return jpeg.Decode(file)
}

func (j *JPEGFormat) Encode(file *os.File, img image.Image) error {
	return jpeg.Encode(file, img, nil)
}

func (p *PNGFormat) Decode(file *os.File) (image.Image, error) {
	return png.Decode(file)
}

func (p *PNGFormat) Encode(file *os.File, img image.Image) error {
	return png.Encode(file, img)
}

func (g *GIFFormat) Decode(file *os.File) (image.Image, error) {
	return gif.Decode(file)
}

func (g *GIFFormat) Encode(file *os.File, img image.Image) error {
	return gif.Encode(file, img, &gif.Options{})
}
