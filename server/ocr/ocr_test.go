package ocr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testScreenshot = "./testdata/test.jpg"
	testDesc       = "A lightweight tool to upload MacOS screenshots to an S3-compatible provider."
)

func TestGosseractOcr_Describe(t *testing.T) {
	ocr, err := NewOCR()
	assert.NoError(t, err)
	defer ocr.Close()

	desc, err := ocr.Describe("./testdata/test.jpg")

	assert.NoError(t, err)
	assert.Equal(t, testDesc, desc)
}

type OCRMock struct{}

func (o *OCRMock) Describe(path string) (string, error) {
	if path == "expected_path" {
		return "expected_result", nil
	}

	if path == "expected_error_path" {
		return "", fmt.Errorf("expected_error")
	}

	return "unexpected_path", nil
}
func (o *OCRMock) Close() error { return nil }

func TestListenOnImages(t *testing.T) {
	ocr := &OCRMock{}
	in := make(chan string)
	out := make(chan []string)

	go ListenOnImages(ocr, in, out)

	in <- "expected_path"
	result := <-out

	assert.Len(t, result, 2)
	assert.Equal(t, "expected_path", result[0])
	assert.Equal(t, "expected_result", result[1])

	in <- "expected_error_path"
	result = <-out

	assert.Len(t, result, 2)
	assert.Equal(t, "expected_error_path", result[0])
	assert.Equal(t, "", result[1])

	close(in)
}
