package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	pb "github.com/moapis/imageapi/imageapi"
	"github.com/moapis/imageapi/s3"
	"google.golang.org/grpc"
	fakesock "google.golang.org/grpc/test/bufconn"
)

func init() {
	port = 9000
	addr = "localhost"
}

func startMock(mock imageServiceServer) *grpc.Server {
	listener := fakesock.Listen(9)
	newServer := grpc.NewServer()
	pb.RegisterImageServiceServer(newServer, mock)
	go func() {
		newServer.Serve(listener)
	}()
	defer listener.Close()
	return newServer
}

func Test_imageServiceServer_NewImageResize(t *testing.T) {
	resp, e := http.Get("https://images.unsplash.com/photo-1535498730771-e735b998cd64?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2134&q=80")
	if e != nil {
		log.Println(e.Error())
	}
	var b []byte
	if b, e = ioutil.ReadAll(resp.Body); e != nil {
		log.Println(e.Error())
	}
	mockServer := imageServiceServer{}
	server := startMock(mockServer)
	ctx, cf := context.WithTimeout(context.Background(), time.Second*3)
	defer cf()
	request := pb.NewImageRequest{}
	request.Image = append(request.Image, b) // append test image
	response, e := mockServer.NewImageResize(ctx, &request)
	if e != nil {
		t.Error(e.Error())
		t.FailNow()
	}
	if response == nil {
		t.Error("Nil response.", e)
		t.FailNow()
	}
	if len(response.Link) != len(request.Image) {
		t.Errorf("%s, %v", response.String(), e)
	}
	key := strings.Split(response.Link[0], "/")[len(response.Link)-1]
	t.Log(response.Link)
	if e := s3.DeleteFile(s3.S3Client, s3.DefaultBucket, key); e != nil {
		t.Error(e.Error())
	}
	server.Stop()
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func TestGetValidContentTypes(t *testing.T) {
	resp, e := http.Get("https://images.unsplash.com/photo-1535498730771-e735b998cd64?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2134&q=80")
	if e != nil {
		log.Println(e.Error())
	}
	var b []byte
	if b, e = ioutil.ReadAll(resp.Body); e != nil {
		log.Println(e.Error())
	}
	result, invalid := getValidContentTypes([][]byte{b})
	if len(result) != 1 || len(invalid) == 1 {
		t.Fail()
	}
}

func TestHaserr(t *testing.T) {
	var e [9]error
	s := "fake error"
	for k := range e {
		if (k % 2) == 0 {
			e[k] = fmt.Errorf("%s at index %d", s, k)
		}
	}
	if !haserr(e[:]) {
		t.Fail()
	}
}
