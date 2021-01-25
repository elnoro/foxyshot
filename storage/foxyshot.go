package storage

import (
	"context"
	"foxyshot/config"
	"foxyshot/server/grpcapi"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewUploader creates a new uploader depedning on the config passed. If both foxyshot and s3 subconfigs are specifed, prefers foxyshot
func NewUploader(c *config.Config) (Uploader, error) {
	if c.Foxyshot.Address != "" {
		return newFoxishotUploader(c.Foxyshot.Address, c.Foxyshot.Insecure, c.Foxyshot.Token)
	}

	return NewS3Uploader(c.S3), nil
}

func newFoxishotUploader(address string, insecure bool, token string) (Uploader, error) {
	opts := []grpc.DialOption{}
	if insecure {
		log.Println("Connecting to an insecure server, make sure you know what you are doing")
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, err
	}
	grpcClient := grpcapi.NewScreenshotsClient(conn)

	return &FoxyshotStorage{client: grpcClient, token: token}, nil
}

// FoxyshotStorage uploads screnshots to an instance of foxyshot server
type FoxyshotStorage struct {
	client grpcapi.ScreenshotsClient
	token  string
}

// Upload a screenshot via an rpc call
func (fs *FoxyshotStorage) Upload(ctx context.Context, path string, options *UploadOptions) (string, error) {
	auth := metadata.AppendToOutgoingContext(ctx, "authorization", fs.token)

	log.Printf("uploading %s to the foxyshot storage\n", path)
	ext := "jpg"
	image, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	screenshot := &grpcapi.Screenshot{}
	screenshot.Extension = ext
	screenshot.Image = image

	l, err := fs.client.Upload(auth, screenshot)
	if err != nil {
		return "", err
	}

	return l.GetLink(), nil
}
