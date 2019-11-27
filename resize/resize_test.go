package resizer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
			bufResult, s, e := ResizeMem(buf, 800, 600)
			if e != nil && !v.ShouldFail {
				t.Errorf("%T | %s", e, e.Error())
			} else if v.ShouldFail {
				t.Log(e.Error())
				t.Skip("ShouldFail:", v.ShouldFail)
			}
			if bufResult == nil || s != strings.Split(ct, "/")[1] {
				t.Error("result buf nil or mime missmatch")
				t.FailNow()
			}
			rb, e := ioutil.ReadAll(bufResult)
			if e != nil {
				t.Error(e.Error())
				t.FailNow()
			}
			if http.DetectContentType(rb) != "png" && !v.ShouldFail {
				t.Error("Content type missmatch.")
				t.FailNow()
			}
		})
	}
}
