package storage_test

import (
	"context"
	"foxyshot/config"
	"foxyshot/storage"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
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

	uploader := newS3Uploader(endpoint)
	url, err := uploader.Upload(ctx, f.Name())
	assert.NoError(t, err)

	assert.Contains(t, url, endpoint+"/expected-bucket/")

	resp, err := http.Get(url)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, 200)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, string(body), uploadContent)
}

func newS3Uploader(endpoint string) storage.Uploader {
	s3Config := &config.S3Config{
		Key:        testUser,
		Secret:     testPass,
		Endpoint:   endpoint,
		Region:     "eu-west-1",
		Bucket:     testBucket,
		PublicURIs: false,
		Duration:   1,
	}
	uploader := storage.NewS3Uploader(s3Config)
	return uploader
}

func startMinio(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "quay.io/minio/minio:latest",
		ExposedPorts: []string{"9000/tcp"},
		WaitingFor:   wait.ForLog("Documentation"), // TODO use http strategy
		Entrypoint:   []string{"/bin/bash", "-c", "mkdir -p /data/" + testBucket + "; minio server /data --console-address \":9001\""},
		Env: map[string]string{
			"MINIO_ROOT_USER":     testUser,
			"MINIO_ROOT_PASSWORD": testPass,
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
