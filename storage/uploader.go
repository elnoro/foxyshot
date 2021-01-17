package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

const (
	defaultBucket   = "foxy"
	defaultDuration = time.Hour
)

// UploadOptions specify details of uploading to files
type UploadOptions struct {
	Bucket string
	// Sets an expiration date for presigned url (only used is ShortPublicUrls is set to false)
	Accessible time.Duration
	// Whether screenshots are made public or private with pre-signed urls
	ShortPublicUrls bool
}

// DefaultOptions are the default upload options for Uploader
var DefaultOptions = &UploadOptions{
	Bucket:          defaultBucket,
	Accessible:      defaultDuration,
	ShortPublicUrls: true,
}

// Uploader Abstract interface for uploading screenshots, other packages should not care if its s3 or gs or whatever
type Uploader interface {
	Upload(ctx context.Context, path string, options *UploadOptions) (string, error)
}

type s3CompatibleUploader struct {
	client *s3.S3
}

func newS3Client(key string, secret string, endpoint string, region string) *s3.S3 {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	return s3Client

}

// NewDigitalOceanUploader creates new Uploader instances compatible with Digital Ocean Spaces API
func NewDigitalOceanUploader(key string, secret string, endpoint string, region string) Uploader {
	c := newS3Client(key, secret, endpoint, region)

	return &s3CompatibleUploader{client: c}
}

// Upload uploads file to s3 and returns presigned url
func (u *s3CompatibleUploader) Upload(ctx context.Context, path string, options *UploadOptions) (string, error) {
	key, err := u.uploadFile(ctx, path, options)
	if err != nil {
		return "", err
	}
	log.Printf("Uploaded %s as %s \n", path, key)

	url, err := u.generateURL(key, options)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (u *s3CompatibleUploader) generateURL(key string, options *UploadOptions) (string, error) {
	if options.ShortPublicUrls {
		url := u.client.Endpoint + "/" + options.Bucket + "/" + key

		return url, nil
	}

	return u.signURL(key, options)
}

func (u *s3CompatibleUploader) uploadFile(ctx context.Context, path string, options *UploadOptions) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	key := generateObjectKey()

	var acl string
	if options.ShortPublicUrls {
		acl = "public-read"
	} else {
		acl = "private"
	}

	input := s3.PutObjectInput{
		Bucket:      aws.String(options.Bucket),
		Key:         aws.String(key),
		Body:        file,
		ACL:         aws.String(acl),
		ContentType: aws.String("image/jpeg"),
	}
	output, err := u.client.PutObjectWithContext(ctx, &input)
	if err != nil {
		return "", err
	}

	log.Printf("Uploaded %s, got %v \n", path, output)

	return key, nil
}

func (u *s3CompatibleUploader) signURL(key string, options *UploadOptions) (string, error) {
	req, _ := u.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(options.Bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(options.Accessible)
	if err != nil {
		return "", err
	}

	return url, nil
}

func generateObjectKey() string {
	timePrefix := time.Now().Format("2006_01_02") // help humans to differentiate between screenshots
	uuid, err := uuid.NewRandom()                 // adding uuid to avoid enumeration
	if err != nil {
		log.Printf("Failed to generate uuid, got error %s\n", err)
	}

	return fmt.Sprintf("%s_%s.jpg", timePrefix, uuid)
}
