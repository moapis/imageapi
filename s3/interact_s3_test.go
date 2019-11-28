package s3

import (
	"bytes"
	"errors"
	minio "github.com/minio/minio-go"
	rs "github.com/moapis/imageapi/resize"
	"io"
	"os"
	"testing"
)

type FakeS3Client struct {
	WantError bool
}

func (c *FakeS3Client) PutObject(string, string, io.Reader, int64, minio.PutObjectOptions) (int64, error) {
	if c.WantError {
		return 0, errors.New("Error test")
	}
	return 1, nil
}
func (c *FakeS3Client) GetObject(string, string, minio.GetObjectOptions) (*minio.Object, error) {
	if c.WantError {
		return new(minio.Object), errors.New("Error test")
	}
	return &minio.Object{}, nil
}
func (c *FakeS3Client) MakeBucket(string, string) error {
	if c.WantError {
		return errors.New("Error test")
	}
	return nil
}
func (c *FakeS3Client) BucketExists(string) (bool, error) {
	if c.WantError {
		return false, errors.New("Error test")
	}
	return true, nil
}
func (c *FakeS3Client) RemoveObject(string, string) error {
	if c.WantError {
		return errors.New("Error test")
	}
	return nil
}

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
	c, e := s3Init()
	if e != nil {
		t.Error(e.Error())
		t.Fail()
	}
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
	s3, e := s3Init()
	if e != nil {
		t.Error(e.Error())
		t.Fail()
	}
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

func Test_putFile(t *testing.T) {
	type args struct {
		c           s3Client
		bucket      string
		keyName     string
		src         bytes.Buffer
		contentType string
	}
	b := &bytes.Buffer{} // mock data object
	b.Write([]byte("some data"))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "empty bucket",
			args: args{c: &FakeS3Client{WantError: true},
				bucket:      "",
				keyName:     "testkey",
				src:         *b,
				contentType: "image/png"}, wantErr: true},
		{name: "empty keyName",
			args: args{c: &FakeS3Client{WantError: true},
				bucket:      "test",
				keyName:     "",
				src:         *b,
				contentType: "image/png"}, wantErr: true},
		{name: "empty buffer",
			args: args{c: &FakeS3Client{WantError: true},
				bucket:      "test",
				keyName:     "testkey",
				src:         bytes.Buffer{},
				contentType: "image/png"}, wantErr: true},
		{name: "actual PUT fails.",
			args: args{c: &FakeS3Client{WantError: true},
				bucket:      "test",
				keyName:     "testkey",
				src:         *b,
				contentType: "image/png"}, wantErr: true},
		{name: "success",
			args: args{c: &FakeS3Client{WantError: false},
				bucket:      "test",
				keyName:     "testkey",
				src:         *b,
				contentType: "image/png"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := putFile(tt.args.c, tt.args.bucket, tt.args.keyName, tt.args.src, tt.args.contentType); (err != nil) != tt.wantErr {
				t.Errorf("putFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getFile(t *testing.T) {
	type args struct {
		c       s3Client
		bucket  string
		keyName string
	}
	t1 := &bytes.Buffer{}
	t1.Write([]byte("data1"))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Get fake",
			args: args{c: &FakeS3Client{WantError: true},
				bucket: DefaultBucket, keyName: "keyName"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getFile(tt.args.c, tt.args.bucket, tt.args.keyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_S3Put(t *testing.T) {
	S3Client = &FakeS3Client{}
	h := &SetterGetter{}
	b := &bytes.Buffer{}
	b.WriteString("not empty")
	if e := h.S3Put(DefaultBucket, "someKey", b, "image/png"); e != nil {
		t.Error(e.Error())
		t.FailNow()
	}
}

// Return object cannot be faked.
func Test_S3Get(t *testing.T) {
	SetCreds()
	var e error
	S3Client, e = s3Init()
	if e != nil {
		t.Error(e.Error())
		t.FailNow()
	}
	key := string(rs.MakeRandomString(24))
	testBucket := string(rs.MakeRandomString(24))
	if e = S3Client.MakeBucket(testBucket, ""); e != nil {
		t.Error(e.Error())
		t.FailNow()
	}
	src := &bytes.Buffer{}
	src.WriteString("test object")
	i, e := S3Client.PutObject(testBucket, key, src, int64(src.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if e != nil {
		t.Error(e.Error(), i)
		t.FailNow()
	}
	h := SetterGetter{}
	b, e := h.S3Get(testBucket, key)
	if e != nil {
		t.Errorf("%+v | %+v", e, b)
		t.FailNow()
	}
	S3Client.RemoveObject(testBucket, key)
}
