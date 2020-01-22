### Image API
- Requires Envoy proxy for usage with browser clients.

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
