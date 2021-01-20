package imageprocessing

import (
	"log"
	"os"
)

func newRemoverPipeline(pipeline ScreenshotPipeline) *RemoverPipeline {
	return &RemoverPipeline{
		pipeline: pipeline,
		remover:  &osRemover{},
	}
}

// RemoverPipeline uses another pipeline to process image and then removes the original image
type RemoverPipeline struct {
	pipeline ScreenshotPipeline
	remover  remover
}

// Run internal pipeline first and then try to remove the original file
func (p *RemoverPipeline) Run(path string) (string, error) {
	processed, err := p.pipeline.Run(path)
	if err != nil {
		return "", err
	}
	go p.remover.Remove(path)

	return processed, nil
}

type remover interface {
	Remove(path string)
}

type osRemover struct {
}

// Remove files with os.Remove
func (r *osRemover) Remove(path string) {
	if err := os.Remove(path); err != nil {
		log.Printf("Could not remove %s, got %v", path, err)
	}
}
