package imageprocessing

import (
	"fmt"
	"foxyshot/config"
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var invalidPngData = []struct {
	file   string
	errMsg string
}{
	{"testdata/notanimage", "png error, unexpected EOF"},
	{"doesnotexist", "png error, open doesnotexist: no such file or directory"},
}

func TestNewPipeline(t *testing.T) {
	c := &config.Config{}

	jpg := NewPipeline(c)
	assert.IsType(t, &readerOptimizer{}, jpg)

	c.Screenshots.RemoveOriginals = true

	remove := NewPipeline(c)
	assert.IsType(t, &RemoverPipeline{}, remove)
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
		_ = os.Remove(f)
	}

	assert.Equal(t, f, "")
	assert.Error(t, err)
	assert.Regexp(t, "no such file or directory$", err.Error())
}

func TestJpgOptimizer_Optimize(t *testing.T) {
	testOptimizer := &jpegOptimizer{tmpFolder: "testdata", quality: 99}

	f, err := testOptimizer.Optimize(mockImage)
	defer func(name string) { _ = os.Remove(name) }(f)

	assert.FileExists(t, f)
	assert.NoError(t, err)
}

func TestJpgOptimizer_OptimizeError(t *testing.T) {
	largeRect := image.Rect(0, 0, 1<<16, 1<<16)
	largeImage := image.NewGray(largeRect)

	testOptimizer := &jpegOptimizer{tmpFolder: "testdata", quality: 99}

	f, err := testOptimizer.Optimize(largeImage)

	assert.Empty(t, f)
	assert.EqualError(t, err, "jpeg optimization error, jpeg: image is too large to encode")
}

func TestNewJpgPipeline(t *testing.T) {
	p := newJpgPipeline(-99)

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

// TODO add benches with larger files to the repo
var benches = []struct {
	name      string
	inputFile string
	quality   int
}{
	{"jpg + smallest png", "testdata/valid.png", 30},
	{"jpg + smallest png", "testdata/valid.png", 100},
	{"jpg + smallest png", "testdata/valid.png", 255},
}

func BenchmarkScreenshot(b *testing.B) {
	for _, set := range benches {
		screenshotPipeline := newJpgPipeline(set.quality)
		b.Run(fmt.Sprintf("%s - quality %d", set.name, set.quality), func(b *testing.B) {
			p, err := screenshotPipeline.Run(set.inputFile)
			if err != nil {
				b.FailNow()
			}
			err = os.Remove(p)
			if err != nil {
				b.FailNow()
			}
		})
	}
}
