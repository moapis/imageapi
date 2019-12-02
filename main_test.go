package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	pb "github.com/moapis/imageapi/imageapi"
	"github.com/moapis/imageapi/models"
	s3 "github.com/moapis/imageapi/s3"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"google.golang.org/grpc"
	fakesock "google.golang.org/grpc/test/bufconn"
)

func init() {
	port = 9000
	addr = "localhost"
}

type FakeSetterGetter struct {
	WantError bool
}

func (sg *FakeSetterGetter) S3Put(bucket string, key string, b *bytes.Buffer, otype string) error {
	if sg.WantError {
		return errors.New("Test error")
	}
	return nil
}
func (sg *FakeSetterGetter) S3Get(bucket string, key string) (*bytes.Buffer, error) {
	if sg.WantError {
		return nil, errors.New("Test error")
	}
	return &bytes.Buffer{}, nil
}

func (*FakeSetterGetter) S3Remove(bucket string, key string) error {
	return nil
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

var fakeByteStringImage = []byte("fakeimage")

func Test_imageServiceServer_NewImageResize_manual(t *testing.T) {
	var b []byte
	var e error
	if b, e = ioutil.ReadFile("resize/test_images/original/original.jpeg"); e != nil {
		log.Println(e.Error())
	}
	mockServer := imageServiceServer{}
	mockServer.S3 = &FakeSetterGetter{}
	db, e = sql.Open("postgres", psqlConnectionURL)
	if e != nil {
		t.Error(e.Error(), psqlConnectionURL)
		t.FailNow()
	}
	server := startMock(mockServer)
	ctx := context.TODO()
	request := pb.NewImageRequest{}
	request.Image = append(request.Image, b) // append test image
	response, e := mockServer.NewImageResize(ctx, &request)
	if e != nil {
		t.Error(e.Error(), psqlConnectionURL)
		t.FailNow()
	}
	if response == nil {
		t.Error("Nil response.", e)
		t.FailNow()
	}
	rsp, e := db.Exec("delete from images where id=$1;", response.Structure[0].GetResizedID())
	if e != nil {
		t.Error(e.Error())
	}
	n, e := rsp.RowsAffected()
	if e != nil {
		t.Error(e.Error())
	}
	if n != 1 {
		t.Errorf("Expected clean-up rows affected %d", 1)
	}
	key := strings.Split(response.Link[0], "/")[len(response.Link)-1]
	t.Log(response.Link)
	s3.S3Client.RemoveObject(s3.DefaultBucket, key)
	server.Stop()
}

func TestGetValidContentTypes_manual(t *testing.T) {
	var b []byte
	var e error
	if b, e = ioutil.ReadFile("resize/test_images/original/original.jpeg"); e != nil {
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

func Test_getValidContentTypes(t *testing.T) {
	type args struct {
		grpcImageSlice [][]byte
	}
	tests := []struct {
		name  string
		args  args
		want  [][]byte
		want1 []int
	}{
		{
			name:  "invalid",
			args:  args{grpcImageSlice: [][]byte{fakeByteStringImage}},
			want1: []int{0}, // this should report that index 0 is invalid
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getValidContentTypes(tt.args.grpcImageSlice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValidContentTypes() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getValidContentTypes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_imageServiceServer_NewImageResize(t *testing.T) {
	type fields struct {
		UnimplementedImageServiceServer pb.UnimplementedImageServiceServer
		S3                              s3.ObjectSetterGetter
	}
	type args struct {
		ctx    context.Context
		images *pb.NewImageRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.NewImageResponse
		wantErr bool
	}{
		{
			name:    "want error",
			fields:  fields{S3: &FakeSetterGetter{WantError: true}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: true,
			want:    &pb.NewImageResponse{},
			args:    args{ctx: context.Background(), images: &pb.NewImageRequest{}},
		},
		{
			name:    "want success",
			fields:  fields{S3: &FakeSetterGetter{WantError: false}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: false,
			want:    &pb.NewImageResponse{},
			args:    args{ctx: context.Background(), images: &pb.NewImageRequest{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := imageServiceServer{
				UnimplementedImageServiceServer: tt.fields.UnimplementedImageServiceServer,
				S3:                              tt.fields.S3,
			}
			got, err := is.NewImageResize(tt.args.ctx, tt.args.images)
			if !tt.wantErr && err != nil {
				t.Errorf("imageServiceServer.NewImageResize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageServiceServer.NewImageResize() = [%v] want [%v]", got, tt.want)
			}
		})
	}
}

func Test_listen(t *testing.T) {
	tests := []struct {
		name string
		want *grpc.Server
	}{
		{name: "server start test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := listen()
			if got == nil {
				t.Errorf("listen() = %v, want %v", got, tt.want)
			}
			got.Stop()
		})
	}
}

func Test_imageServiceServer_NewImagePreserve(t *testing.T) {
	type fields struct {
		UnimplementedImageServiceServer pb.UnimplementedImageServiceServer
		S3                              s3.ObjectSetterGetter
	}
	type args struct {
		ctx    context.Context
		images *pb.NewImageRequest
	}
	b1, _ := ioutil.ReadFile("resize/test_images/original/original.jpg")
	b1Arr := [][]byte{b1}
	b2, _ := ioutil.ReadFile("resize/test_images/original/original.fakeext")
	b2Arr := [][]byte{b2}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.NewImageResponse
		wantErr bool
	}{
		{
			name:    "want error",
			fields:  fields{S3: &FakeSetterGetter{WantError: true}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: true,
			want:    &pb.NewImageResponse{},
			args:    args{ctx: context.Background(), images: &pb.NewImageRequest{Image: b2Arr}},
		},
		{
			name:    "want success",
			fields:  fields{S3: &FakeSetterGetter{WantError: false}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: false,
			want:    &pb.NewImageResponse{},
			args:    args{ctx: context.Background(), images: &pb.NewImageRequest{Image: b1Arr}},
		},
	}
	var e error
	db, e = sql.Open("postgres", psqlConnectionURL)
	if e != nil {
		t.Error(e.Error(), psqlConnectionURL)
		t.FailNow()
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := imageServiceServer{
				UnimplementedImageServiceServer: tt.fields.UnimplementedImageServiceServer,
				S3:                              tt.fields.S3,
			}
			got, err := is.NewImagePreserve(tt.args.ctx, tt.args.images)
			if err != nil && !tt.wantErr {
				t.Errorf("imageServiceServer.NewImagePreserve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.Contains(got.GetLink()[0], fmt.Sprintf("%s/%s", s3.S3Endpoint, s3.DefaultBucket)) {
				t.Error("Link doesn't match.")
				t.FailNow()
			}
			if !tt.wantErr {
				rsp, e := db.Exec("delete from images where id=$1;", got.Structure[0].GetOriginalID())
				if e != nil {
					t.Error(e.Error())
				}
				n, e := rsp.RowsAffected()
				if e != nil {
					t.Error(e.Error())
				}
				if n != 1 {
					t.Errorf("Expected clean-up rows affected %d | %+v", 1, got)
				}
			}
		})
	}
}

func Test_imageServiceServer_RemoveImage(t *testing.T) {
	type fields struct {
		UnimplementedImageServiceServer pb.UnimplementedImageServiceServer
		S3                              s3.ObjectSetterGetter
	}
	type args struct {
		ctx     context.Context
		request *pb.RemoveImageRequest
	}
	link := "somelink"
	var e error
	db, e = sql.Open("postgres", psqlConnectionURL)
	if e != nil {
		t.Error(e.Error(), psqlConnectionURL)
		t.FailNow()
	}
	m := models.Image{LinkOriginal: null.NewString(link, true)}
	m.Insert(context.TODO(), db, boil.Infer())
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.RemoveImageResponse
		wantErr bool
	}{
		{
			name:    "want error",
			fields:  fields{S3: &FakeSetterGetter{WantError: true}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: true,
			want:    &pb.RemoveImageResponse{},
			args:    args{ctx: context.Background(), request: &pb.RemoveImageRequest{Link: []string{"    ", "fakelink"}}},
		},
		{
			name:    "success",
			fields:  fields{S3: &FakeSetterGetter{WantError: false}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: false,
			want:    &pb.RemoveImageResponse{Status: "OK"},
			args:    args{ctx: context.Background(), request: &pb.RemoveImageRequest{Link: []string{link}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := imageServiceServer{
				UnimplementedImageServiceServer: tt.fields.UnimplementedImageServiceServer,
				S3:                              tt.fields.S3,
			}
			got, err := is.RemoveImage(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageServiceServer.RemoveImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageServiceServer.RemoveImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageServiceServer_NewImageResizeAndPreserve(t *testing.T) {
	type fields struct {
		UnimplementedImageServiceServer pb.UnimplementedImageServiceServer
		S3                              s3.ObjectSetterGetter
	}
	b1, _ := ioutil.ReadFile("resize/test_images/original/original.jpg")
	b1Arr := [][]byte{b1}
	b2, _ := ioutil.ReadFile("resize/test_images/original/original.fakeext")
	b2Arr := [][]byte{b2}
	var e error
	db, e = sql.Open("postgres", psqlConnectionURL)
	if e != nil {
		t.Error(e.Error(), psqlConnectionURL)
		t.FailNow()
	}
	type args struct {
		ctx    context.Context
		images *pb.NewImageRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "want success",
			fields:  fields{S3: &FakeSetterGetter{WantError: false}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: false,
			args:    args{ctx: context.Background(), images: &pb.NewImageRequest{Image: b1Arr}},
		},
		{
			name:    "want error",
			fields:  fields{S3: &FakeSetterGetter{WantError: true}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: true,
			args:    args{ctx: context.Background(), images: &pb.NewImageRequest{Image: b2Arr}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := imageServiceServer{
				UnimplementedImageServiceServer: tt.fields.UnimplementedImageServiceServer,
				S3:                              tt.fields.S3,
			}
			got, err := is.NewImageResizeAndPreserve(tt.args.ctx, tt.args.images)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageServiceServer.NewImageResizeAndPreserve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			response := got.GetStructure()
			if !tt.wantErr && len(response) > 0 && (!strings.Contains(response[0].GetOriginalLink(), s3.DefaultBucket) || !strings.Contains(response[0].GetResizedLink(), s3.DefaultBucket)) {
				t.Errorf("Expected valid links received %+v", got)
			}
			if !tt.wantErr && response == nil {
				t.Errorf("Expected valid links received %+v", got)
			}
			if !tt.wantErr { // clean inserts
				rt, e := db.Exec("delete from images where id=$1;", response[0].GetOriginalID())
				if e != nil {
					t.Error(e.Error())
				}
				ra, _ := rt.RowsAffected()
				if ra != 1 {
					t.Errorf("Expected db insert clean-up RowsAffected %d got %d", 1, ra)
				}
			}
		})
	}
}

func Test_imageServiceServer_NewImageResizeAtDimensions(t *testing.T) {
	type fields struct {
		UnimplementedImageServiceServer pb.UnimplementedImageServiceServer
		S3                              s3.ObjectSetterGetter
	}
	b1, _ := ioutil.ReadFile("resize/test_images/original/original.jpg")
	b1Arr := [][]byte{b1}
	b2, _ := ioutil.ReadFile("resize/test_images/original/original.fakeext")
	b2Arr := [][]byte{b2}
	var e error
	db, e = sql.Open("postgres", psqlConnectionURL)
	if e != nil {
		t.Error(e.Error(), psqlConnectionURL)
		t.FailNow()
	}
	type args struct {
		ctx    context.Context
		images *pb.NewImageRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "want s3 error",
			fields:  fields{S3: &FakeSetterGetter{WantError: true}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: true,
			args:    args{ctx: context.Background(), images: &pb.NewImageRequest{Image: b1Arr, Dimensions: &pb.ImageDimensions{Width: 300, Height: 300}}},
		},
		{
			name:    "want success",
			fields:  fields{S3: &FakeSetterGetter{WantError: false}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: false,
			args:    args{ctx: context.Background(), images: &pb.NewImageRequest{Image: b1Arr, Dimensions: &pb.ImageDimensions{Width: 300, Height: 300}}},
		},
		{
			name:    "invalid dimensions",
			fields:  fields{S3: &FakeSetterGetter{WantError: false}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: true,
			args:    args{ctx: context.Background(), images: &pb.NewImageRequest{Image: b1Arr}},
		},
		{
			name:    "want error not image",
			fields:  fields{S3: &FakeSetterGetter{WantError: true}, UnimplementedImageServiceServer: pb.UnimplementedImageServiceServer{}},
			wantErr: true,
			args:    args{ctx: context.Background(), images: &pb.NewImageRequest{Image: b2Arr, Dimensions: &pb.ImageDimensions{Width: 300, Height: 300}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := imageServiceServer{
				UnimplementedImageServiceServer: tt.fields.UnimplementedImageServiceServer,
				S3:                              tt.fields.S3,
			}
			_, err := is.NewImageResizeAtDimensions(tt.args.ctx, tt.args.images)
			if (err == nil) && tt.wantErr {
				t.Errorf("imageServiceServer.NewImageResizeAtDimensions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
