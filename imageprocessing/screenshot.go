package imageprocessing

import (
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

const (
	// DefaultPrefix for optimized screenshot files
	DefaultPrefix = "foxy_img"
	// DefaultTmpFolder where optimized screenshot files will be stored
	DefaultTmpFolder = "/tmp"
)

// ScreenshotPipeline is an interface to optimization pipeline for images
// Currently only PNG - JPG pipeline is supported
type ScreenshotPipeline interface {
	// Accepts path to an existing image and returns path to an optimized image
	Run(path string) (string, error)
}

type readerOptimizer struct {
	reader    screenshotReader
	optimizer screenshotOptimizer
}

func (pipeline *readerOptimizer) Run(path string) (string, error) {
	img, err := pipeline.reader.Read(path)
	if err != nil {
		return "", err
	}

	return pipeline.optimizer.Optimize(img)
}

// NewJpgPipeline Creates ScreenshotPipeline that converts images into jpgs
// for MacOS quality 30 seems to be sufficient for screenshots and provides up to 90% savings in file size
func NewJpgPipeline(quality int) ScreenshotPipeline {
	jpegOptimizer := &jpegOptimizer{
		quality:   quality,
		tmpFolder: DefaultTmpFolder,
		prefix:    DefaultPrefix,
	}
	pngReader := &pngReader{}

	return &readerOptimizer{reader: pngReader, optimizer: jpegOptimizer}
}

// screenshotOptimizer is an interface for optimizing screenshot images
type screenshotOptimizer interface {
	// Optimize returns path to a temporary file containing optimized image
	//It is the caller's responsibility to remove the file when no longer needed.
	Optimize(img image.Image) (string, error)
}

// jpegOptimizer saves image to jpg with a specified quality
type jpegOptimizer struct {
	tmpFolder string
	prefix    string
	quality   int
}

func (opt *jpegOptimizer) Optimize(img image.Image) (string, error) {
	file, err := ioutil.TempFile(opt.tmpFolder, opt.prefix)
	if err != nil {
		return "", err
	}
	defer file.Close()

	log.Println("Saving compressed screenshot ", file.Name())

	err = jpeg.Encode(file, img, &jpeg.Options{Quality: opt.quality})
	if err != nil {
		os.Remove(file.Name())
		return "", err
	}

	return file.Name(), nil
}

// screenshotReader is an interface for reading screenshots into an image.Image
type screenshotReader interface {
	Read(path string) (image.Image, error)
}

type pngReader struct{}

func (reader *pngReader) Read(path string) (image.Image, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	// TODO check if using image.Decode makes sense - which format do Monosnap, Joxi etc. use?
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
