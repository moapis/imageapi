package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"

	pb "github.com/moapis/imageapi/imageapi"
	rs "github.com/moapis/imageapi/resize"
	s3 "github.com/moapis/imageapi/s3"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// IMAGEAPI_PORT (env var)
var port int

// IMAGEAPI_ADDR (env var)
var addr string

const (
	defaultWidth  = 800
	defaultHeight = 600
)

type imageServiceServer struct {
	pb.UnimplementedImageServiceServer
	S3 s3.ObjectSetterGetter
}

const tmpStore = "/tmp/image_api_data"

var validMimeTypes = [4]string{"jpg", "jpeg", "png", "gif"}

func checkMime(data []byte) bool {
	for _, tp := range validMimeTypes {
		if fmt.Sprintf("image/%s", tp) == http.DetectContentType(data) {
			return true
		}
	}
	return false
}

// scanMime is supposed to get the bad indexes out of the way so that they can be skipped on the slower loops.
func scanMime(b [][]byte) []int {
	var invalidContentTypeIndexes []int
	for k, data := range b {
		if !checkMime(data) {
			invalidContentTypeIndexes = append(invalidContentTypeIndexes, k)
		}
	}
	return invalidContentTypeIndexes
}

func getValidContentTypes(grpcImageSlice [][]byte) ([][]byte, []int) {
	var dataArray [][]byte
	invalidIndexes := scanMime(grpcImageSlice)
	for k, v := range grpcImageSlice {
		invalid := false
		for _, ii := range invalidIndexes {
			if k == ii {
				invalid = true
			}
		}
		if !invalid {
			dataArray = append(dataArray, v)
		}
	}
	return dataArray, invalidIndexes
}

// Log every error encountered.
func haserr(err []error) bool {
	has := false
	for _, e := range err {
		if e != nil {
			log.Println(e.Error())
			has = true
		}
	}
	return has
}

// Uploads and resizes a variable size of images keeping and returning only the resized version.
func (is imageServiceServer) NewImageResize(ctx context.Context, images *pb.NewImageRequest) (*pb.NewImageResponse, error) {
	dataArray, invalidArray := getValidContentTypes(images.GetImage())
	response := pb.NewImageResponse{}
	for _, b := range dataArray {
		buf := new(bytes.Buffer)
		var e [2]error
		_, e[0] = buf.Write(b)
		var s string
		buf, s, e[1] = rs.ResizeMem(buf, defaultWidth, defaultHeight)
		if haserr(e[:]) {
			return &response, status.Error(codes.InvalidArgument, "Internal server error, check server logs for aditional information.")
		}
		key := string(rs.MakeRandomString(12))
		if e := is.S3.S3Put(s3.DefaultBucket, key, buf, fmt.Sprintf("image/%s", s)); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.InvalidArgument, "Internal server error, check server logs for aditional information.")
		}
		response.Link = append(response.Link, fmt.Sprintf("https://%s/%s/%s", s3.S3Endpoint, s3.DefaultBucket, key))
	}
	if len(invalidArray) > 0 {
		return &response, status.Error(codes.InvalidArgument, fmt.Sprintf("Invalid content type found at indexes: %+v", invalidArray))
	}
	return &response, nil
}

//Uploads image while not attempting to alter it.
func (is imageServiceServer) NewImagePreserve(ctx context.Context, images *pb.NewImageRequest) (response *pb.NewImageResponse, e error) {
	return
}

// Uploads image and keeps both versions of the file.
func (is imageServiceServer) NewImageResizeAndPreserve(ctx context.Context, images *pb.NewImageRequest) (response *pb.NewImageResponse, e error) {
	return
}

// Resizes image at specified dimensions
func (is imageServiceServer) NewImageResizeAtDimensions(ctx context.Context, images *pb.NewImageRequest) (response *pb.NewImageResponse, e error) {
	return
}

func init() {
	info, _ := debug.ReadBuildInfo()
	log.Printf("%+v", info.Main.Path)
	var e error
	port, e = strconv.Atoi(os.Getenv("IMAGEAPI_PORT"))
	if e != nil {
		log.Println(e.Error())
	}
	addr = os.Getenv("IMAGEAPI_ADDR")
}

func listen() *grpc.Server {
	listener, e := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if e != nil {
		log.Println(e.Error())
	}
	is := imageServiceServer{}
	is.S3 = &s3.SetterGetter{} // Inject actual implementation.
	s := grpc.NewServer()
	pb.RegisterImageServiceServer(s, is)
	go func() {
		s.Serve(listener)
	}()
	return s
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	log.Printf("Listening on %s:%d ...", addr, port)
	grpcServer := listen()
	<-sig
	log.Println("Stopping server.")
	grpcServer.GracefulStop()
}
