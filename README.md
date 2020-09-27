### Image API

Image API is a microservice that integrates and is reliant on https://github.com/moapis/authenticator for validating the tokens that it receives from the client and it is a specialized service that handles image processing and storage to S3 compatible storages directly from buffer without writing to disk. 
It interacts with a PostgreSQL database for storing relevant information. 
It can handle not only single/batch upload and image storage but also image resizing and overlaying one image over another and can also store mp4 videos. 
- This is a RESTful API that is meant to run as a microservice in Docker along with Authenticator and Shop.
- Client requests are made to this API using Web gRPC protocol and require Envoy proxy to translate the request into native gRPC.

Accepted mime types: image/jpg, image/jpeg, image/png, image/gif, image/bmp, video/mp4

It uses https://github.com/minio/minio-go client library for interacting with S3.

It also uses https://github.com/anthonynsimon/bild for some of the image manipulation.

For more information regarding Google's gRPC protocol please visit: https://grpc.io/about/

Envoy docs: https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/other_protocols/grpc

To regenerate the Go code from image_api.proto file please consult:
https://grpc.io/docs/languages/go/basics/#generating-client-and-server-code

To generate web gRPC client Javascript code please refer to:
https://grpc.io/blog/grpc-web-ga/

-----------------------------------------------------------------------------------------------------------------------------


#### Configuration
Below are the various environment variables that need to be propagated when deploying this service.


Postgresql database environment variables.

- IMAGEAPI_PQ_HOST
- IMAGEAPI_PQ_PORT
- IMAGEAPI_PQ_USER
- IMAGEAPI_PQ_PASS
- IMAGEAPI_PQ_DBNAME
- IMAGEAPI_PQ_SSLMODE

Image API grpc server configuration.

- IMAGEAPI_ADDR
- IMAGEAPI_PORT

Authenticator endpoint that it should exchange tokens with.

- AUTHENTICATOR_HOSTNAME

Exposed S3 storage client environment variables.

- S3_ENDPOINT
- S3_KEY
- S3_SECRET
- S3_TLS
- S3_DEFAULT_BUCKET

Below env vars default to minio if not set.

- S3_TEST_ENDPOINT
- S3_TEST_KEY
- S3_TEST_SECRET
- S3_TEST_TLS
- S3_TEST_BUCKET
