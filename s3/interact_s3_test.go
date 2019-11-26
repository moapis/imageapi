package s3

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/minio/minio-go"
	rs "github.com/moapis/imageapi/resize"
)

var testImage = "https://picsum.photos/id/192/1920/1080"
var testImageType = "jpeg"

func SetCreds() {
	if os.Getenv("S3_TEST_ENDPOINT") == "" {
		S3Endpoint = "play.minio.io:9000"
	}
	if os.Getenv("S3_TEST_BUCKET") == "" {
		DefaultBucket = "magick-crop"
	}
	if os.Getenv("S3_TEST_KEY") == "" {
		S3Key = "Q3AM3UQ867SPQQA43P2F"
	}
	if os.Getenv("S3_TEST_SECRET") == "" {
		S3Secret = "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	}
	if os.Getenv("S3_TEST_TLS") == "" {
		tls = true
	}
}

func TestS3Init(t *testing.T) {
	SetCreds()
	c := s3Init()
	if o, e := c.BucketExists(DefaultBucket); e != nil || !o {
		t.Errorf("S3Init test failed: %v, %t", o, e)
		if e != nil {
			t.Errorf("S3Init test failed: %+v",
				e.(minio.ErrorResponse))
		}
	}
	if o, e := c.BucketExists("fake2374876376234"); e != nil {
		t.Errorf("S3Init test failed: %+v",
			e.(minio.ErrorResponse))
	} else if o {
		t.Errorf("Bucket exists returns false positive.")
	}
}
func TestMakeBucket(t *testing.T) {
	SetCreds()
	bkt := string(rs.MakeRandomString(15))
	s3 := s3Init()
	if e := s3.MakeBucket(bkt, ""); e != nil {
		t.Errorf("Failed at bucket creation test: %s", e.Error())
	}
	if o, e := s3.BucketExists(string(bkt)); e != nil || !o {
		t.Errorf("S3Init test failed: %s", e.Error())
	}
	if e := s3.MakeBucket("#$$#$@!", ""); e == nil {
		t.FailNow()
	} else {
		switch e := e.(type) {
		case minio.ErrorResponse:
			t.Logf("%+v", e)
			if e.StatusCode != 403 {
				t.Errorf("%+v", e)
				t.FailNow()
			}
		default:
			if e.Error() != "Bucket name contains invalid characters" {
				t.Fail()
			}
		}
	}
}
func TestPutFile(t *testing.T) {
	SetCreds()
	s3 := s3Init()
	response, e := http.Get(testImage)
	if e != nil {
		t.Error("Fail[1]:", e.Error())
	}
	b, e := ioutil.ReadAll(response.Body)
	if e != nil {
		t.Error("Fail[2]:", e.Error())
	}
	fileName := string(rs.MakeRandomString(15))
	folder := os.TempDir()
	saveFilePath := strings.Join([]string{folder, "/", fileName, ".", testImageType}, "")
	if e = ioutil.WriteFile(saveFilePath, b, os.ModePerm); e != nil {
		t.Error("Fail[3]:", e.Error())
	}
	if e = PutFile(s3, saveFilePath, DefaultBucket, fileName, fmt.Sprintf("image/%s", testImageType)); e != nil {
		t.Error("Fail[4]:", http.DetectContentType(b), saveFilePath, e.Error())
	}
	url := strings.Join([]string{"https:/", S3Endpoint, DefaultBucket, fileName}, "/")
	response2, e := http.Get(url)
	if e != nil {
		t.Error("Fail[5]:", e.Error())
	}
	b2, e := ioutil.ReadAll(response2.Body)
	if e != nil {
		t.Error("Fail[6]:", e.Error())
	}
	imgType := http.DetectContentType(b2)
	if imgType != fmt.Sprintf("image/%s", testImageType) {
		t.Errorf("S3 PUT test failed at image type check: %s", imgType)
	}
	os.Remove(saveFilePath)
	if e := DeleteFile(s3, DefaultBucket, S3Key); e != nil {
		t.Error(e.Error())
	}
	// test fail
	if e := PutFile(s3, "/tmp/fakepath", DefaultBucket, "fakekey39280", "image/jpg"); e == nil {
		t.Error("Expected error at s3 putFile()")
	} else {
		switch e := e.(type) {
		case minio.ErrorResponse:
			if e.StatusCode != 404 {
				t.Error(e.Error())
			}
		case *os.PathError:
			t.Log(e.Error())
		}
	}
}

func TestDeleteFile(t *testing.T) {
	// test delete error
	SetCreds()
	s3 := s3Init()
	if e := DeleteFile(s3, "fakebucket328497", "fake33423243847892337248"); e == nil {
		t.Error("Expected error from deleteFile()\n")
	} else {
		t.Logf("%+v", e.(minio.ErrorResponse))
		if e.(minio.ErrorResponse).Code != "NoSuchBucket" {
			t.Fail()
		}
	}
}

func TestGetFile(t *testing.T) {
	SetCreds()
	bucket := DefaultBucket
	path := os.TempDir()
	key := "get_test"
	response, _ := http.Get(testImage)
	b, _ := ioutil.ReadAll(response.Body)
	filepath := path + "/" + key
	if err := ioutil.WriteFile(filepath, b, os.ModePerm); err != nil {
		t.Error(err.Error())
	}
	c := s3Init()
	if err := PutFile(c, filepath, bucket, key, fmt.Sprintf("image/%s", testImageType)); err != nil {
		log.Println(err.Error())
	}
	if err := getFile(c, bucket, key, filepath); err != nil {
		t.Error(err.Error())
	}
	fl, err := os.Open(filepath)
	if err != nil || fl == nil {
		t.Error(err.Error())
	}
	falseKey := "false_key"
	fpath := fmt.Sprintf("/tmp/%s", falseKey)
	if err := getFile(c, bucket, falseKey, fpath); err == nil {
		t.Error("No error when expected.")
	} else {
		t.Log(err.Error())
	}
	if e := DeleteFile(c, bucket, key); e != nil {
		t.Error(e.Error())
	}
}

func TestConstructURL(t *testing.T) {
	SetCreds()
	bkt := "fake"
	test := fmt.Sprintf("https://%s.%s", bkt, S3Endpoint)
	test1 := fmt.Sprintf("http://%s.%s", bkt, S3Endpoint)
	if !strings.Contains(constructURL(bkt), test) {
		t.Errorf("Expected %s but got %s", test, constructURL(bkt))
	}
	tls = false
	if !strings.Contains(constructURL(bkt), test1) {
		t.Errorf("Expected %s but got %s", test1, constructURL(bkt))
	}
}

func TestGenerateURL(t *testing.T) {
	SetCreds()
	c := s3Init()
	path := os.TempDir()
	key := "generate_url_test"
	response, _ := http.Get(testImage)
	b, _ := ioutil.ReadAll(response.Body)
	filepath := path + "/" + key
	if err := ioutil.WriteFile(filepath, b, os.ModePerm); err != nil {
		t.Error(err.Error())
	}
	if e := PutFile(c, filepath, DefaultBucket, key, testImageType); e != nil {
		t.Error(e.Error())
	}
	url, e := generateURL(c, DefaultBucket, key)
	if e != nil {
		log.Println(e.Error())
	}
	expected := fmt.Sprintf("/%s/%s", DefaultBucket, key)
	if url.Path == "" {
		t.FailNow()
	} else if url.Path != expected {
		t.Logf("Expected %s but got %s", expected, url.Path)
		t.Fail()
	}
}
