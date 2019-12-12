FROM golang:1.13 as b
RUN mkdir app
WORKDIR /app
COPY . .
ENV GO111MODULE=on
RUN go build -o imageapi_binary
######################################
FROM debian:stretch-slim
COPY --from=b /app/imageapi_binary .

#ENV IMAGEAPI_ADDR=0.0.0.0
#ENV IMAGEAPI_PORT=8081

#ENV S3_ENDPOINT=play.minio.io:9000
#ENV S3_KEY=Q3AM3UQ867SPQQA43P2F
#ENV S3_SECRET=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
#ENV S3_DEFAULT_BUCKET=magick-crop

#ENV IMAGEAPI_PQ_HOST=some-postgres
#ENV IMAGEAPI_PQ_PORT=5432
#ENV IMAGEAPI_PQ_USER=postgres
#ENV IMAGEAPI_PQ_PASS=mysecretpassword
#ENV IMAGEAPI_PQ_DBNAME=imageapi
#ENV IMAGEAPI_PQ_SSLMODE=disable
#ENV AUTHENTICATOR_HOSTNAME=
RUN apt-get update 
RUN apt-get install -y ca-certificates

ENTRYPOINT ./imageapi_binary
