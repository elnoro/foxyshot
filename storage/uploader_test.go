package storage_test

import (
	"context"
	"fmt"
	"foxyshot/config"
	"foxyshot/storage"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testUser      = "FOXYSHOT_USER"
	testPass      = "FOXYSHOT_PASS"
	testBucket    = "expected-bucket"
	uploadContent = "successfully uploaded"
)

func Test_UploadHappyPath(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.Background()
	minioC, err := startMinio(ctx)
	assert.NoError(t, err)
	defer func(minioC testcontainers.Container, ctx context.Context) {
		err := minioC.Terminate(ctx)
		if err != nil {
			log.Println("Container termination failed. Please clean up manually!")
		}
	}(minioC, ctx)

	endpoint, err := minioC.Endpoint(ctx, "http")
	assert.NoError(t, err)

	f, err := createUploadFile(uploadContent)
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	tests := []struct {
		name       string
		publicURIs bool
	}{
		{
			"public uris",
			true,
		},
		{
			"presigned urls",
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uploader := newS3Uploader(endpoint, test.publicURIs)

			url, err := uploader.Upload(ctx, f.Name())
			fmt.Println(url)
			assert.NoError(t, err)

			assert.Contains(t, url, fmt.Sprintf("%s/%s/", endpoint, testBucket))

			resp, err := http.Get(url)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, string(body), uploadContent)
		})
	}
}

func newS3Uploader(endpoint string, publicURIs bool) storage.Uploader {
	s3Config := &config.S3Config{
		Key:      testUser,
		Secret:   testPass,
		Region:   "eu-west-1",
		Bucket:   testBucket,
		Duration: 60 * time.Second,

		Endpoint:   endpoint,
		PublicURIs: publicURIs,
	}
	uploader := storage.NewS3Uploader(s3Config)
	return uploader
}

func startMinio(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "bitnami/minio",
		ExposedPorts: []string{"9000/tcp"},
		WaitingFor:   wait.ForHTTP("/" + testBucket).WithMethod("HEAD"),
		Env: map[string]string{
			"MINIO_ROOT_USER":       testUser,
			"MINIO_ROOT_PASSWORD":   testPass,
			"MINIO_DEFAULT_BUCKETS": testBucket + ":public",
		},
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func createUploadFile(contents string) (*os.File, error) {
	f, err := os.CreateTemp("", "upload_test")
	if err != nil {
		return nil, err
	}
	_, err = f.WriteString(contents)
	return f, err
}
