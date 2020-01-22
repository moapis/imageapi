package resizer

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestResize(t *testing.T) {
	pwd, e := os.Getwd()
	if e != nil {
		t.Errorf("%T | %s", e, e.Error())
		t.FailNow()
	}
	ext := "jpeg"
	b, e := ioutil.ReadFile(fmt.Sprintf("%s/test_images/original/original.%s", pwd, ext))
	originalCt := http.DetectContentType(b)
	initialLength := len(b) // inital length
	filename := string(MakeRandomString(15))
	tmp := os.TempDir()
	testFile := strings.Join([]string{tmp, "/", filename, ".", ext}, "")
	e = ioutil.WriteFile(testFile, b, os.ModePerm)
	if e != nil {
		t.Error(e.Error())
	}
	if e = Resize(testFile, 150, 150); e != nil {
		t.Error(e.Error())
	}
	b, e = ioutil.ReadFile(testFile)
	resultCt := http.DetectContentType(b)
	if initialLength == len(b) {
		t.Error("Fail: expected smaller byte slice.", "Initial:", initialLength, ", resized:", len(b))
	}
	if originalCt != resultCt {
		t.Error("Mime type missmatch.")
	}
	fakepath := "/fakepath/to/foo"
	if e = Resize(fakepath, 150, 150); e == nil {
		t.Error("Expected error using path:", fakepath)
	}
}
func TestResizeMem(t *testing.T) {
	cases := []struct {
		Extension  string
		ShouldFail bool
	}{
		{Extension: "jpg", ShouldFail: false},
		{Extension: "jpeg", ShouldFail: false},
		{Extension: "png", ShouldFail: false},
		{Extension: "gif", ShouldFail: false},
		{Extension: "bmp", ShouldFail: false},
		{Extension: "fakeext", ShouldFail: true},
	}
	for k, v := range cases {
		t.Run(fmt.Sprintf("%d # Testing %s", k, v.Extension), func(t *testing.T) {
			path, e := os.Getwd()
			if e != nil {
				t.Errorf("%T | %s", e, e.Error())
				t.FailNow()
			}
			b, e := ioutil.ReadFile(fmt.Sprintf("%s/test_images/original/original.%s", path, v.Extension))
			if e != nil {
				t.Error(e.Error(), path)
				t.FailNow()
			}
			ct := http.DetectContentType(b)
			buf := new(bytes.Buffer)
			if n, e := buf.Write(b); e != nil || n == 0 {
				t.Error(e.Error())
				t.FailNow()
			}
			s, e := ResizeMem(buf, 800, 600)
			if e != nil && !v.ShouldFail {
				t.Errorf("%T | %s", e, e.Error())
			} else if v.ShouldFail {
				t.Log(e.Error())
				t.Skip("ShouldFail:", v.ShouldFail)
			}
			if buf == nil || s != strings.Split(ct, "/")[1] {
				t.Error("result buf nil or mime missmatch")
				t.FailNow()
			}
			rb, e := ioutil.ReadAll(buf)
			if e != nil {
				t.Error(e.Error())
				t.FailNow()
			}
			if http.DetectContentType(rb) != ct && !v.ShouldFail {
				t.Error("Content type missmatch.")
				t.FailNow()
			}
		})
	}
}
func getImgBuffers(img1, img2 string) (*bytes.Buffer, *bytes.Buffer, error) {
	ilogo, e := ioutil.ReadFile(img1)
	if e != nil {
		return nil, nil, e
	}
	ipng, e := ioutil.ReadFile(img2)
	if e != nil {
		return nil, nil, e
	}
	blogo := new(bytes.Buffer)
	bpng := new(bytes.Buffer)
	_, e = blogo.Write(ilogo)
	if e != nil {
		return nil, nil, e
	}
	_, e = bpng.Write(ipng)
	if e != nil {
		return nil, nil, e
	}
	return blogo, bpng, nil
}
func TestOverlay(t *testing.T) {
	blogo, bpng, e := getImgBuffers("test_images/original/original.png", "test_images/original/original.png")
	if e != nil {
		t.Error(e.Error())
	}
	blogo1, bpng1, e := getImgBuffers("test_images/original/original.jpeg", "test_images/original/original.jpeg")
	if e != nil {
		t.Error(e.Error())
	}
	blogo2, bpng2, e := getImgBuffers("test_images/original/original.gif", "test_images/original/original.gif")
	if e != nil {
		t.Error(e.Error())
	}
	blogo3, bpng3, e := getImgBuffers("test_images/original/original.gif", "test_images/original/original.bmp")
	if e != nil {
		t.Error(e.Error())
	}
	blogo4, bpng4, e := getImgBuffers("test_images/original/original.gif", "test_images/original/original.bmp")
	if e != nil {
		t.Error(e.Error())
	}
	blogo5, bpng5, e := getImgBuffers("test_images/original/original.gif", "test_images/original/original.bmp")
	if e != nil {
		t.Error(e.Error())
	}
	type args struct {
		img1     *bytes.Buffer
		img2     *bytes.Buffer
		position string
		rX       int
		rY       int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test #1",
			args:    args{blogo, bpng, "bottomright", 100, 100},
			wantErr: false,
		},
		{name: "test #2",
			args:    args{blogo1, bpng1, "bottomleft", 100, 100},
			wantErr: false,
		},
		{name: "test #3",
			args:    args{blogo2, bpng2, "topright", 100, 100},
			wantErr: false,
		},
		{name: "test #4",
			args:    args{img1: blogo3, img2: bpng3, position: "topleft", rX: 100, rY: 100},
			wantErr: false,
		},
		{name: "test #5",
			args:    args{img1: blogo4, img2: bpng4, position: "center", rX: 700, rY: 700},
			wantErr: false,
		},
		{name: "test #6",
			args:    args{img1: blogo5, img2: new(bytes.Buffer), position: "topleft", rX: 100, rY: 100},
			wantErr: true,
		},
		{name: "test #7",
			args:    args{img1: new(bytes.Buffer), img2: bpng5, position: "topleft", rX: 100, rY: 100},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			itype, err := Overlay(tt.args.img1, tt.args.img2, tt.args.position, tt.args.rX, tt.args.rY)
			if (err != nil) && !tt.wantErr {
				t.Errorf("Overlay() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				b, e := ioutil.ReadAll(tt.args.img2)
				if e != nil {
					t.Error(e.Error())
					t.FailNow()
				}
				if e = ioutil.WriteFile(fmt.Sprintf("%s.%s", "test_images/overlay_test", itype), b, os.ModePerm); e != nil {
					t.Error(e.Error())
				}
			}
		})
	}
}

func Test_calcOverlayPosition(t *testing.T) {
	type args struct {
		rectangleunder image.Rectangle
		rectangleover  image.Rectangle
		position       string
	}
	tests := []struct {
		name string
		args args
		want image.Point
	}{
		{name: "test bottomright",
			want: image.Point{-400, -400},
			args: args{image.Rect(0, 0, 1200, 900), image.Rect(0, 0, 800, 500), "bottomright"},
		},
		{name: "test bottomleft",
			want: image.Point{0, -400},
			args: args{image.Rect(0, 0, 1200, 900), image.Rect(0, 0, 800, 500), "bottomleft"},
		},
		{name: "test bottomright",
			want: image.Point{-400, 0},
			args: args{image.Rect(0, 0, 1200, 900), image.Rect(0, 0, 800, 500), "topright"},
		},
		{name: "test bottomright",
			want: image.Point{-200, -200},
			args: args{image.Rect(0, 0, 1200, 900), image.Rect(0, 0, 800, 500), "center"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcOverlayPosition(tt.args.rectangleunder, tt.args.rectangleover, tt.args.position); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calcOverlayPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}
