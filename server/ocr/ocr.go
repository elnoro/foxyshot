package ocr

import (
	"log"

	gosseract "github.com/otiai10/gosseract/v2"
)

// OCR Interface to generate descriptions for the images
type OCR interface {
	Describe(path string) (string, error)
	Close() error
}

type gosseractOcr struct {
	cl *gosseract.Client
}

func (o *gosseractOcr) Describe(path string) (out string, err error) {
	err = o.cl.SetImage(path)
	if err != nil {
		return out, err
	}
	out, err = o.cl.Text()
	if err != nil {
		return out, err
	}

	return out, err
}

func (o *gosseractOcr) Close() error {
	return o.cl.Close()
}

// NewOCR creates a new OCR instance to parse incoming images
func NewOCR() (OCR, error) {
	cl := gosseract.NewClient()

	return &gosseractOcr{cl: cl}, nil
}

// ListenOnImages waits on any paths to images from a channel and sends ocr results into another channel
// if the channel is closed, the function closes ocr
func ListenOnImages(ocr OCR, paths chan string, pathsAndDescriptions chan []string) {
	for p := range paths {
		desc, err := ocr.Describe(p)
		if err != nil {
			pathsAndDescriptions <- []string{p, ""}
		}
		pathsAndDescriptions <- []string{p, desc}
	}
	if err := ocr.Close(); err != nil {
		log.Println(err)
	}
}
