package resizer

import (
	"bytes"
	"image"
	gif "image/gif"
	jpeg "image/jpeg"
	png "image/png"
	"io"
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
func ResizeMem(r io.Reader, w, h int) (*bytes.Buffer, string, error) {
	img, s, e := image.Decode(r)
	if e != nil {
		return nil, "", e
	}
	log.Printf("Decoded type [%s]", s)
	transform.Resize(img, w, h, transform.Linear)
	buf := new(bytes.Buffer)
	switch s {
	case "png":
		e = png.Encode(buf, img)
	case "jpeg":
		e = jpeg.Encode(buf, img, &jpeg.Options{Quality: 100})
	case "gif":
		e = gif.Encode(buf, img, nil)
	case "bmp":
		e = bmp.Encode(buf, img)
	}
	return buf, s, e
}
