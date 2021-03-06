package resizer

import (
	"bytes"
	"image"
	"image/draw"
	gif "image/gif"
	jpeg "image/jpeg"
	png "image/png"
	"log"
	"math/rand"
	"time"

	bmp "golang.org/x/image/bmp"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

func genPseudoRand() *rand.Rand {
	rInt := rand.New(rand.NewSource(rand.Int63() * time.Now().UnixNano()))
	return rInt
}

// MakeRandomString generates a pseudo random string with the length specified as parameter.
func MakeRandomString(bytesLength int) []byte {
	byteVar := make([]byte, bytesLength)
	chars := "abcdefghijklmnopqrstuvwxyz123456789" // our posibilities
	for i := range byteVar {
		x := genPseudoRand()
		byteVar[i] = chars[x.Intn(len(chars))]
	}
	return byteVar
}

// Resize uses bild library to open convert and write the image to the same path.
func Resize(imagePath string, w, h int) error {
	i, e := imgio.Open(imagePath)
	if e != nil {
		return e
	}
	resized := transform.Resize(i, w, h, transform.Linear)
	e = imgio.Save(imagePath, resized, imgio.JPEGEncoder(100))
	return e
}

// ResizeMem - resize without writing to disk.
// warning: jpg becomes jpeg when this is used.
// defaults to png
func ResizeMem(r *bytes.Buffer, w, h int) (string, error) {
	img, s, e := image.Decode(r)
	if e != nil {
		return "", e
	}
	rect := img.Bounds()
	r.Reset()
	log.Printf("Decoded type [%s]", s)
	oh := float64(rect.Dy())
	ow := float64(rect.Dx())
	var ar float64
	ar = ow / oh
	nh := float64(w) / ar
	imgc := transform.Resize(img, w, int(nh), transform.Linear)
	log.Printf("New resize: [w %d - h %d]", w, int(nh))
	switch s {
	case "png":
		e = png.Encode(r, imgc)
	case "jpeg":
		e = jpeg.Encode(r, imgc, &jpeg.Options{Quality: 100})
	case "gif":
		e = gif.Encode(r, imgc, nil)
	case "bmp":
		e = bmp.Encode(r, imgc)
	}
	return s, e
}
func calcOverlayPosition(rectangleunder, rectangleover image.Rectangle, position string) image.Point {
	var p image.Point
	switch position {
	case "bottomright":
		p.X = rectangleunder.Max.X - rectangleover.Max.X
		p.Y = rectangleunder.Max.Y - rectangleover.Max.Y
	case "bottomleft":
		p.Y = rectangleunder.Max.Y - rectangleover.Max.Y
	case "topright":
		p.X = rectangleunder.Max.X - rectangleover.Max.X
	case "center":
		p.X = (rectangleunder.Max.X - rectangleover.Max.X) / 2
		p.Y = (rectangleunder.Max.Y - rectangleover.Max.Y) / 2
	}
	p.X *= -1
	p.Y *= -1
	return p
}

// Overlay applies a stamp of img1 on img2 at specified coordinates.
// Takes buffers as arguments as to prevent unnecessary writes to disk.
// The overlay image should be smaller to obtain a stamp or watermark effect.
func Overlay(img1, img2 *bytes.Buffer, position string, resizeX, resizeY int) (string, error) {
	if _, e := ResizeMem(img1, resizeX, resizeY); e != nil {
		return "", e
	}
	iover, _, e := image.Decode(img1)
	if e != nil {
		return "", e
	}
	iunder, s, e := image.Decode(img2)
	if e != nil {
		return "", e
	}
	dst := image.NewRGBA(iunder.Bounds())
	draw.Draw(dst, iunder.Bounds(), iunder, image.Point{X: 0, Y: 0}, draw.Src) // draw the background
	draw.Draw(dst, iunder.Bounds(), iover, calcOverlayPosition(iunder.Bounds(), iover.Bounds(), position), draw.Over)
	img2.Reset()
	switch s {
	case "png":
		e = png.Encode(img2, dst)
	case "jpeg":
		e = jpeg.Encode(img2, dst, &jpeg.Options{Quality: 100})
	case "gif":
		e = gif.Encode(img2, dst, nil)
	case "bmp":
		e = bmp.Encode(img2, dst)
	}
	return s, e
}
