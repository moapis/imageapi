package s3

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"

	minio "github.com/minio/minio-go"
)

// S3Endpoint - S3_ENDPOINT
var S3Endpoint string
var tls = true

// S3Key - S3_KEY
var S3Key string

// S3Secret - S3_SECRET
var S3Secret string

// DefaultBucket - S3_DEFAULT_BUCKET
var DefaultBucket string

// Basic Minio s3 client signature.
type s3Client interface {
	PutObject(bucket string, key string, src io.Reader, sz int64, opts minio.PutObjectOptions) (int64, error)
	GetObject(bucket string, key string, opts minio.GetObjectOptions) (*minio.Object, error)
	MakeBucket(bucket string, location string) error
	BucketExists(bucket string) (bool, error)
	RemoveObject(bucket string, key string) error
}

// S3Client - initialized handler using env credentials.
var S3Client s3Client

// ObjectSetterGetter - Interface that provides signature setter/getter from s3.
type ObjectSetterGetter interface {
	S3Put(bucket string, key string, b *bytes.Buffer, otype string) error
	S3Get(bucket string, key string) (*bytes.Buffer, error)
	S3Remove(bucket string, key string) error
}

// SetterGetter - implements ObjectSetterGetter using Minio S3 Client.
type SetterGetter struct{}

// S3Put - Put Object to S3
func (*SetterGetter) S3Put(bucket string, key string, b *bytes.Buffer, otype string) error {
	return putFile(S3Client, bucket, key, *b, otype)
}

// S3Get - Get Object from S3
func (*SetterGetter) S3Get(bucket string, key string) (*bytes.Buffer, error) {
	return getFile(S3Client, bucket, key)
}

// S3Remove - Remove single object from S3
func (*SetterGetter) S3Remove(bucket string, key string) error {
	return S3Client.RemoveObject(bucket, key)
}

func init() {
	S3Endpoint = os.Getenv("S3_ENDPOINT")
	S3Key = os.Getenv("S3_KEY")
	S3Secret = os.Getenv("S3_SECRET")
	DefaultBucket = os.Getenv("S3_DEFAULT_BUCKET")
	var e error
	S3Client, e = s3Init()
	if e != nil {
		log.Fatal(e.Error())
	}
}

func setPolicy(Client *minio.Client, bucket string) {
	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + bucket + `/*"]}]}`
	e := Client.SetBucketPolicy(bucket, policy)
	if e != nil {
		log.Println(e.Error())
	}
}

func s3Init() (s3Client, error) {
	log.Println("S3:", S3Endpoint, "tls:", tls)
	Client, e := minio.New(S3Endpoint, S3Key, S3Secret, tls)
	if e != nil {
		log.Printf("On S3 init encountered: %+v\n", e.(minio.ErrorResponse).Code)
		return nil, e
	}
	return Client, e
}

// putFile - internal - Uploads object to specified bucket using specified keyname. Cleans up the file after successful PUT.
func putFile(c s3Client, bucket string, keyName string, src bytes.Buffer, contentType string) error {
	if bucket == "" {
		return errors.New("bucket is empty string")
	}
	if keyName == "" {
		return errors.New("keyName is empty string")
	}
	if src.Len() == 0 {
		return errors.New("src buffer is empty")
	}
	_, e := c.PutObject(bucket, keyName, &src, int64(src.Len()), minio.PutObjectOptions{ContentType: contentType})
	if e != nil {
		return e
	}
	return nil
}

// getFile - internal - Download object.
func getFile(c s3Client, bucket string, keyName string) (*bytes.Buffer, error) {
	reader, e := c.GetObject(bucket, keyName, minio.GetObjectOptions{})
	if e != nil {
		return nil, e
	}
	defer reader.Close()
	b := &bytes.Buffer{}
	_, e = io.Copy(b, reader)
	return b, e
}
