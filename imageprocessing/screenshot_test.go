package imageprocessing

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var invalidPngData = []struct {
	file   string
	errMsg string
}{
	{"testdata/notanimage", "unexpected EOF"},
	{"doesnotexist", "open doesnotexist: no such file or directory"},
}

func TestPngReader_ReadInvalidData(t *testing.T) {
	testReader := &pngReader{}

	for _, d := range invalidPngData {
		img, err := testReader.Read(d.file)
		assert.Nil(t, img)
		assert.EqualError(t, err, d.errMsg)
	}
}

func TestPngReader_ReadValidData(t *testing.T) {
	testReader := &pngReader{}
	img, err := testReader.Read("testdata/valid.png")

	assert.NotNil(t, img)
	assert.NoError(t, err)
}

func TestJpgOptimizer_OptimizeInaccessibleFolder(t *testing.T) {
	testOptimizer := &jpegOptimizer{tmpFolder: "doesnotexist", quality: 99}

	f, err := testOptimizer.Optimize(mockImage)
	if f != "" {
		os.Remove(f)
	}

	assert.Equal(t, f, "")
	assert.Error(t, err)
	assert.Regexp(t, "no such file or directory$", err.Error())
}

func TestJpgOptimizer_Optimize(t *testing.T) {
	testOptimizer := &jpegOptimizer{tmpFolder: "testdata", quality: 99}

	f, err := testOptimizer.Optimize(mockImage)
	defer os.Remove(f)

	assert.FileExists(t, f)
	assert.NoError(t, err)
}

func TestJpgOptimizer_OptimizeError(t *testing.T) {
	largeRect := image.Rect(0, 0, 1<<16, 1<<16)
	largeImage := image.NewGray(largeRect)

	testOptimizer := &jpegOptimizer{tmpFolder: "testdata", quality: 99}

	f, err := testOptimizer.Optimize(largeImage)

	assert.Empty(t, f)
	assert.EqualError(t, err, "jpeg: image is too large to encode")
}

func TestNewJpgPipeline(t *testing.T) {
	p := NewJpgPipeline(-99)

	assert.IsType(t, &readerOptimizer{}, p)
}

type Mock struct {
}

func (m *Mock) Read(path string) (image.Image, error) {
	if path == "expected path" {
		return mockImage, nil
	}

	return mockImage, fmt.Errorf("read error")
}

func (m *Mock) Optimize(img image.Image) (string, error) {
	if img == mockImage {
		return "expected result", nil
	}

	return "", fmt.Errorf("optimize error")
}

func TestReaderOptimizer_Run(t *testing.T) {
	m := &Mock{}
	ro := &readerOptimizer{reader: m, optimizer: m}

	f, err := ro.Run("expected path")

	assert.Equal(t, "expected result", f)
	assert.NoError(t, err)
}

func TestReaderOptimizer_RunReaderError(t *testing.T) {
	m := &Mock{}
	ro := &readerOptimizer{reader: m, optimizer: m}

	f, err := ro.Run("wrong path")

	assert.Equal(t, "", f)
	assert.EqualError(t, err, "read error")
}

var mockImage = image.NewGray(image.Rect(0, 0, 1, 1))
