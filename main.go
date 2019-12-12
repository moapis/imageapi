package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	authenticator "github.com/moapis/authenticator/pb"
	"github.com/moapis/authenticator/verify"
	"github.com/pascaldekloe/jwt"

	pb "github.com/moapis/imageapi/imageapi"
	"github.com/moapis/imageapi/models"
	rs "github.com/moapis/imageapi/resize"
	s3 "github.com/moapis/imageapi/s3"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var psqlConnectionURL string

var db *sql.DB
var dbuser string
var dbpassword string
var dbname string
var dbhost string
var dbport string
var sslmode string
var authServer = ""

// IMAGEAPI_PORT (env var)
var port int

// IMAGEAPI_ADDR (env var)
var addr string

const (
	defaultWidth                 = 500
	defaultHeight                = 350
	errorStringInternal          = "Internal server error, check server logs for aditional information."
	errorInvalidContentFound     = "Invalid content type found at indexes: %+v"
	errorPermissionDenied        = "Token validation failed."
	invalidDimensionsErrorString = "Invalid resize dimensions supplied."
)

var clientTokenChecker verify.Verificator

type imageServiceServer struct {
	pb.UnimplementedImageServiceServer
	S3             s3.ObjectSetterGetter
	ccTokenChecker *grpc.ClientConn
}

const tmpStore = "/tmp/image_api_data"

var validMimeTypes = [5]string{"jpg", "jpeg", "png", "gif", "bmp"}

func checkMime(data []byte) bool {
	for _, tp := range validMimeTypes {
		if fmt.Sprintf("image/%s", tp) == http.DetectContentType(data) {
			return true
		}
	}
	log.Printf("Detected type: [%s]\n", http.DetectContentType(data))
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
		s, e[1] = rs.ResizeMem(buf, defaultWidth, defaultHeight)
		if haserr(e[:]) {
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		key := string(rs.MakeRandomString(15))
		if e := is.S3.S3Put(s3.DefaultBucket, key, buf, fmt.Sprintf("image/%s", s)); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		link := fmt.Sprintf("https://%s/%s/%s", s3.S3Endpoint, s3.DefaultBucket, key)
		newLink := models.Image{LinkResized: null.NewString(link, true)}
		if e := newLink.Insert(ctx, db, boil.Infer()); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		response.Link = append(response.Link, link)
		response.Structure = append(response.Structure, &pb.NewImageResponseStruct{ResizedLink: link, ResizedID: uint32(newLink.ID)})
	}
	if len(invalidArray) > 0 {
		return &response, status.Error(codes.InvalidArgument, fmt.Sprintf(errorInvalidContentFound, invalidArray))
	}
	return &response, nil
}

//Uploads image while not attempting to alter it.
func (is imageServiceServer) NewImagePreserve(ctx context.Context, images *pb.NewImageRequest) (*pb.NewImageResponse, error) {
	dataArray, invalidArray := getValidContentTypes(images.GetImage())
	response := pb.NewImageResponse{}
	for _, b := range dataArray {
		buf := new(bytes.Buffer)
		_, e := buf.Write(b)
		if e != nil {
			log.Println(e.Error())
		}
		s := http.DetectContentType(b)
		key := string(rs.MakeRandomString(15))
		if e := is.S3.S3Put(s3.DefaultBucket, key, buf, s); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		link := fmt.Sprintf("https://%s/%s/%s", s3.S3Endpoint, s3.DefaultBucket, key)
		newLink := models.Image{LinkOriginal: null.NewString(link, true)}
		if e := newLink.Insert(ctx, db, boil.Infer()); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		response.Link = append(response.Link, link)
		response.Structure = append(response.Structure, &pb.NewImageResponseStruct{OriginalLink: link, OriginalID: uint32(newLink.ID)})
	}
	if len(invalidArray) > 0 {
		return &response, status.Error(codes.InvalidArgument, fmt.Sprintf(errorInvalidContentFound, invalidArray))
	}
	return &response, nil
}

// Uploads image and keeps both versions of the file.
func (is imageServiceServer) NewImageResizeAndPreserve(ctx context.Context, images *pb.NewImageRequest) (*pb.NewImageResponse, error) {
	dataArray, invalidArray := getValidContentTypes(images.GetImage())
	response := pb.NewImageResponse{}
	for _, b := range dataArray {
		bufOriginal := new(bytes.Buffer)
		bufResized := new(bytes.Buffer)
		if _, e := bufOriginal.Write(b); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		if _, e := bufResized.Write(b); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		s, e := rs.ResizeMem(bufResized, defaultWidth, defaultHeight)
		if e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		keyo := string(rs.MakeRandomString(15))
		keyr := fmt.Sprintf("%s_resized", keyo)
		if e = is.S3.S3Put(s3.DefaultBucket, keyo, bufOriginal, s); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		linko := fmt.Sprintf("https://%s/%s/%s", s3.S3Endpoint, s3.DefaultBucket, keyo)
		if e = is.S3.S3Put(s3.DefaultBucket, keyr, bufResized, s); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		linkr := fmt.Sprintf("https://%s/%s/%s", s3.S3Endpoint, s3.DefaultBucket, keyr)
		newRow := models.Image{LinkOriginal: null.NewString(linko, true), LinkResized: null.NewString(linkr, true)}
		if e = newRow.Insert(ctx, db, boil.Infer()); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		response.Structure = append(response.Structure,
			&pb.NewImageResponseStruct{OriginalLink: linko, ResizedLink: linkr, OriginalID: uint32(newRow.ID), ResizedID: uint32(newRow.ID)})
	}
	if len(invalidArray) > 0 {
		return &response, status.Error(codes.InvalidArgument, fmt.Sprintf(errorInvalidContentFound, invalidArray))
	}
	return &response, nil
}

// Resizes image at specified dimensions
func (is imageServiceServer) NewImageResizeAtDimensions(ctx context.Context, images *pb.NewImageRequest) (*pb.NewImageResponse, error) {
	dataArray, invalidArray := getValidContentTypes(images.GetImage())
	response := pb.NewImageResponse{}
	dimensions := images.GetDimensions()
	w := dimensions.GetHeight()
	h := dimensions.GetWidth()
	if w == 0 || h == 0 {
		log.Println("Invalid image dimensions.")
		return &response, status.Error(codes.InvalidArgument, invalidDimensionsErrorString)
	}
	for _, b := range dataArray {
		buf := new(bytes.Buffer)
		_, e := buf.Write(b)
		if e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		var s string
		s, e = rs.ResizeMem(buf, int(w), int(h))
		if e := is.S3.S3Put(s3.DefaultBucket, s3.DefaultBucket, buf, s); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		key := rs.MakeRandomString(15)
		link := fmt.Sprintf("https://%s/%s/%s", s3.S3Endpoint, s3.DefaultBucket, key)
		mdl := models.Image{LinkResized: null.NewString(link, true)}
		if e = mdl.Insert(ctx, db, boil.Infer()); e != nil {
			log.Println(e.Error())
			return &response, status.Error(codes.Internal, errorStringInternal)
		}
		response.Structure = append(response.Structure, &pb.NewImageResponseStruct{ResizedLink: link, ResizedID: uint32(mdl.ID)})
	}
	if len(invalidArray) > 0 {
		return &response, status.Error(codes.InvalidArgument, fmt.Sprintf(errorInvalidContentFound, invalidArray))
	}
	return &response, nil
}
func (is imageServiceServer) RemoveImage(ctx context.Context, request *pb.RemoveImageRequest) (*pb.RemoveImageResponse, error) {
	links := request.GetLink()
	var e error
	if len(links) == 0 {
		return &pb.RemoveImageResponse{}, status.Error(codes.InvalidArgument, errorStringInternal)
	}
	for _, v := range links {
		link := strings.TrimSpace(v)
		if link == "" {
			return &pb.RemoveImageResponse{}, status.Error(codes.InvalidArgument, errorStringInternal)
		}
		var n int64
		n, e = models.Images(qm.Where("link_original=?", link), qm.Or("link_resized=?", link)).DeleteAll(ctx, db)
		if e != nil {
			log.Println(e.Error())
			return &pb.RemoveImageResponse{}, status.Error(codes.Internal, errorStringInternal)
		}
		if n == 0 {
			log.Println("Expected positive AffectedRows.")
			return &pb.RemoveImageResponse{}, status.Error(codes.Internal, errorStringInternal)
		}
		sl := strings.Split(link, "/")
		if e = is.S3.S3Remove(s3.DefaultBucket, sl[len(sl)-1]); e != nil {
			log.Println(e.Error())
			return &pb.RemoveImageResponse{}, status.Error(codes.Internal, errorStringInternal)
		}
	}
	return &pb.RemoveImageResponse{Status: "OK"}, nil
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
	authServer = os.Getenv("AUTHENTICATOR_HOSTNAME")
	psqlConnectionURL = strings.Join([]string{"postgres://", os.Getenv("IMAGEAPI_PQ_USER"), ":", os.Getenv("IMAGEAPI_PQ_PASS"),
		"@", os.Getenv("IMAGEAPI_PQ_HOST"), ":", os.Getenv("IMAGEAPI_PQ_PORT"), "/", os.Getenv("IMAGEAPI_PQ_DBNAME"), "?sslmode=", os.Getenv("IMAGEAPI_PQ_SSLMODE")}, "")
}

func (is *imageServiceServer) getClientListener(remote string) error {
	var e error
	is.ccTokenChecker, e = grpc.Dial(remote, grpc.WithInsecure())
	return e
}
func (is *imageServiceServer) tokenCheckClientInit() {
	a := authenticator.NewAuthenticatorClient(is.ccTokenChecker)
	clientTokenChecker = verify.Verificator{Client: a}
}

func isNewImageRq(rq interface{}) (*pb.NewImageRequest, bool) {
	r, ok := rq.(*pb.NewImageRequest)
	return r, ok
}
func isRemoveImageRq(rq interface{}) (*pb.RemoveImageRequest, bool) {
	r, ok := rq.(*pb.RemoveImageRequest)
	return r, ok
}
func verifyTkn(ctx context.Context, tkn string) (*jwt.Claims, error) {
	jwtClaims, e := clientTokenChecker.Token(ctx, tkn)
	if e != nil {
		return nil, e
	}
	return jwtClaims, nil
}

// Applicable for all gRPC service methods.
func (is *imageServiceServer) tokenCheckInterceptor(ctx context.Context, rq interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var jwtc *jwt.Claims
	var e error
	if request, ok := isNewImageRq(rq); ok {
		jwtc, e = verifyTkn(ctx, request.GetTkn())
	}
	if request, ok := isRemoveImageRq(rq); ok {
		jwtc, e = verifyTkn(ctx, request.GetTkn())
	}
	if (info != nil) && (jwtc != nil) {
		log.Printf("Method: [%s]\n Issued: [%s] Subject: [%s]  \nerror: [%v]\n", info.FullMethod, jwtc.Issued.String(), jwtc.Subject, e)
	}
	if e != nil {
		return nil, status.Error(codes.PermissionDenied, errorPermissionDenied)
	}
	return handler(ctx, rq)
}

func listen() *grpc.Server {
	listener, e := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if e != nil {
		log.Println(e.Error())
		os.Exit(1)
	}
	is := imageServiceServer{}
	is.S3 = &s3.SetterGetter{} // Inject actual implementation.
	e = is.getClientListener(authServer)
	if e != nil {
		log.Println(e.Error())
		os.Exit(1)
	}
	is.tokenCheckClientInit()
	s := grpc.NewServer(grpc.UnaryInterceptor(is.tokenCheckInterceptor))
	pb.RegisterImageServiceServer(s, is)
	go func() {
		s.Serve(listener)
	}()
	return s
}

func main() {
	var e error
	db, e = sql.Open("postgres", psqlConnectionURL)
	if e != nil {
		log.Println(e.Error())
		os.Exit(1)
	}
	log.Println(psqlConnectionURL)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	log.Printf("Listening on %s:%d ...", addr, port)
	grpcServer := listen()
	<-sig
	log.Println("Stopping server.")
	grpcServer.GracefulStop()
}
