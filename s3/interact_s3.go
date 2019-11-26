package s3

import (
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	minio "github.com/minio/minio-go"
)

// S3Endpoint - S3_ENDPOINT
var S3Endpoint string
var tls bool

// S3Key - S3_KEY
var S3Key string

// S3Secret - S3_SECRET
var S3Secret string

// DefaultBucket - S3_DEFAULT_BUCKET
var DefaultBucket string

// S3Client - initialized handler using env credentials.
var S3Client *minio.Client

func init() {
	S3Endpoint = os.Getenv("S3_ENDPOINT")
	S3Key = os.Getenv("S3_KEY")
	S3Secret = os.Getenv("S3_SECRET")
	s3tls := os.Getenv("S3_TLS")
	if s3tls != "" {
		o, e := strconv.ParseBool(s3tls)
		if e != nil {
			log.Println(e.Error())
		}
		tls = o
	} else {
		log.Println("S3_TLS env not specified.")
	}
	DefaultBucket = os.Getenv("S3_DEFAULT_BUCKET")
	S3Client = s3Init()
}

// For public buckets
func constructURL(bucket string) string {
	var url strings.Builder
	if tls {
		url.WriteString("https://")
	} else {
		url.WriteString("http://")
	}
	url.WriteString(strings.Join([]string{bucket, ".", S3Endpoint}, ""))
	return url.String()
}

func generateURL(c *minio.Client, bucket string, object string) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"xo.jpg\"")
	return c.PresignedGetObject(bucket, object, time.Second*60*60*24*7, reqParams)
}

func setPolicy(Client *minio.Client, bucket string) {
	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + bucket + `/*"]}]}`
	e := Client.SetBucketPolicy(bucket, policy)
	if e != nil {
		log.Println(e.Error())
	}
}

func s3Init() *minio.Client {
	log.Println("S3:", S3Endpoint, "tls:", tls)
	Client, e := minio.New(S3Endpoint, S3Key, S3Secret, tls)
	if e != nil {
		log.Printf("On S3 init encountered: %+v\n", e.(minio.ErrorResponse).Code)
	}
	return Client
}

func makeBucket(Client *minio.Client, bucket string) {
	location := ""
	err := Client.MakeBucket(bucket, location)
	if err != nil {
		log.Println(err.(minio.ErrorResponse).Message, err.(minio.ErrorResponse).StatusCode)
		exists, err := Client.BucketExists(bucket)
		if err != nil {
			log.Println(err.(minio.ErrorResponse).Message, err.(minio.ErrorResponse).StatusCode)
		} else {
			if exists {
				log.Println("Bucket exists.")
			}
		}
	} else {
		log.Println("Bucket", bucket, "created.")
	}
}

// policy.json
func getPolicy(Client *minio.Client, bucket string) error {
	policy, e := Client.GetBucketPolicy(bucket)
	if e != nil {
		log.Println(e.Error())
		return e
	}
	policyReader := strings.NewReader(policy)
	fl, e := os.Create("policy.json")
	if e != nil {
		return e
	}
	n, e := io.Copy(fl, policyReader)
	if e != nil {
		return e
	}
	log.Printf("policy.json created, written %d bytes to disk", n)
	return nil
}

// PutFile - Uploads object to specified bucket using specified keyname. Cleans up the file after successful PUT.
func PutFile(c *minio.Client, filePath string, bucket string, keyName string, contentType string) error {
	fileObject, e := os.Open(filePath)
	if e != nil {
		return e
	}
	fileInfo, e := fileObject.Stat()
	if e != nil {
		return e
	}
	n, e := c.PutObject(bucket, keyName, fileObject, fileInfo.Size(), minio.PutObjectOptions{ContentType: contentType})
	if e != nil {
		return e
	}
	log.Println("Uploaded file of size ", n)
	fileObject.Close()
	return os.Remove(filePath)
}

func DeleteFile(c *minio.Client, bucket, name string) error {
	if e := c.RemoveObject(bucket, name); e != nil {
		return e
	}
	return nil
}

func getFile(c *minio.Client, bucket string, keyName string, newLocalFile string) error {
	reader, e := c.GetObject(bucket, keyName, minio.GetObjectOptions{})
	if e != nil {
		return e
	}
	defer reader.Close()
	newLocFile, _ := os.Create(newLocalFile)
	defer newLocFile.Close()
	info, e := reader.Stat()
	if e != nil {
		log.Printf("%s | %d | %+v\n",
			e.(minio.ErrorResponse).Message,
			e.(minio.ErrorResponse).StatusCode,
			info.Metadata)
		return e
	}
	_, e = io.CopyN(newLocFile, reader, info.Size)
	if e != nil {
		return e
	}
	return nil
}
