package resizer

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestTmpSaveFile(t *testing.T) {
	response, e := http.Get("https://via.placeholder.com/1500")
	if e != nil {
		t.Error(e.Error())
	}
	var b2 []byte
	b, _ := ioutil.ReadAll(response.Body)
	location, e := TmpSaveFile(b)
	if e != nil {
		t.Fail()
	}
	if b2, e = ioutil.ReadFile(location); e != nil || len(b) == 0 {
		t.Error("Fail: ", e.Error())
	}
	if len(b) != len(b2) {
		t.Error(len(b), len(b2))
	}
}

func TestResize(t *testing.T) {
	response, e := http.Get("https://via.placeholder.com/1500")
	if e != nil {
		t.Error(e.Error())
	}
	b, _ := ioutil.ReadAll(response.Body)
	initialLength := len(b) // inital length
	filename := string(MakeRandomString(15))
	path, _ := os.UserHomeDir()
	testFile := strings.Join([]string{path, "/", filename, ".png"}, "")
	e = ioutil.WriteFile(testFile, b, os.ModePerm)
	if e != nil {
		t.Error(e.Error())
	}
	Resize(testFile, 150, 150)
	b, e = ioutil.ReadFile(testFile)
	if initialLength <= len(b) {
		t.Error("Fail: expected smaller byte slice.", "Initial:", initialLength, ", resized:", len(b))
	}
}
func TestResizeMem(t *testing.T) {
	testImg := "https://via.placeholder.com/1500"
	expectedType := "png"
	testImg2 := "https://picsum.photos/id/192/1920/1080"
	expectedType2 := "jpeg"
	response, e := http.Get(testImg)
	// # 1
	if e != nil {
		t.Error(e.Error())
	}
	buf, s, e := ResizeMem(response.Body, 800, 600)
	if e != nil {
		t.Errorf("%T | %s", e, e.Error())
	}
	if buf.Len() == 0 || s != expectedType {
		t.Fail()
	}
	// # 2
	response2, e := http.Get(testImg2)
	if e != nil {
		t.Error(e.Error())
	}
	buf2, s, e := ResizeMem(response2.Body, 800, 600)
	if e != nil {
		t.Errorf("%T | %s", e, e.Error())
	}
	if buf2.Len() == 0 || s != expectedType2 {
		t.Fail()
	}
}
