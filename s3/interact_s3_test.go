package s3

import (
	"fmt"
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
