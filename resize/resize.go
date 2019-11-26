package resizer

import (
	"bytes"
	"image"
	gif "image/gif"
	jpeg "image/jpeg"
	png "image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// TmpSaveFile attempts to create tmp folder in $HOME and save the file with random name.
func TmpSaveFile(b []byte) (string, error) {
	var ext string
	switch http.DetectContentType(b) {
	case "image/jpg":
		ext = ".jpg"
	case "image/jpeg":
		ext = ".jpeg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	}
	path, e := os.UserHomeDir()
	if e != nil {
		return "", e
	}
	folder := filepath.Join(path, "tmp")
	os.Mkdir(folder, os.ModePerm)
	location := filepath.Join(folder, strings.Join([]string{string(MakeRandomString(30)), ext}, ""))
	if e = ioutil.WriteFile(location, b, os.ModePerm); e != nil {
		return "", e
	}
	log.Println(location)
	return location, nil
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
	case "jpg":
		e = jpeg.Encode(buf, img, &jpeg.Options{Quality: 100})
	case "gif":
		e = gif.Encode(buf, img, nil)
	default:
		e = png.Encode(buf, img)
	}
	return buf, s, e
}
