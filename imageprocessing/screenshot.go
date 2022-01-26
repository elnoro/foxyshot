package imageprocessing

import (
	"fmt"
	"foxyshot/config"
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

// NewPipeline uses config to construct the pipeline
func NewPipeline(c *config.Config) ScreenshotPipeline {
	p := newJpgPipeline(c.JpegQuality)
	if c.RemoveOriginals {
		return newRemoverPipeline(p)
	}

	return p
}

// ScreenshotPipeline is an interface to optimization pipeline for images
// Currently only PNG - JPG pipeline is supported
type ScreenshotPipeline interface {
	// Run accepts path to an existing image and returns path to an optimized image
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

// newJpgPipeline Creates ScreenshotPipeline that converts images into jpgs
// for MacOS quality 30 seems to be sufficient for screenshots and provides up to 90% savings in file size
func newJpgPipeline(quality int) ScreenshotPipeline {
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
	verbose   bool
}

func (opt *jpegOptimizer) Optimize(img image.Image) (string, error) {
	file, err := ioutil.TempFile(opt.tmpFolder, opt.prefix)
	if err != nil {
		return "", fmt.Errorf("jpeg error, %w", err)
	}

	if opt.verbose {
		log.Println("Saving compressed screenshot ", file.Name())
	}

	err = jpeg.Encode(file, img, &jpeg.Options{Quality: opt.quality})
	if err != nil {
		rerr := os.Remove(file.Name())
		if rerr != nil {
			return "", fmt.Errorf("invalid jpeg removal error %v, original reason %w", rerr, err)
		}

		return "", fmt.Errorf("jpeg optimization error, %w", err)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf("jpeg error, %w", err)
	}

	return file.Name(), nil
}

// screenshotReader is an interface for reading screenshots into an image.Image
type screenshotReader interface {
	Read(path string) (image.Image, error)
}

type pngReader struct{}

func (reader *pngReader) Read(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("png error, %w", err)
	}
	defer func(file *os.File) {
		cerr := file.Close()
		if cerr != nil {
			err = fmt.Errorf("png error, cannot close file - %w", cerr) // passing error to the top
		}
	}(file)

	// TODO check if using image.Decode makes sense - which format do Monosnap, Joxi etc. use?
	img, err = png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("png error, %w", err)
	}

	return img, nil
}
