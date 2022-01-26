package imageprocessing

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsRemover_Remove(t *testing.T) {
	r := &osRemover{}

	f, err := os.CreateTemp("testdata/", "tempfile")
	if err != nil {
		assert.FailNow(t, "Cannot run this test - no file was created")
	}

	r.Remove(f.Name())

	assert.False(t, assert.FileExists(new(testing.T), f.Name()), "Temp file not removed by remover")

	err = f.Close()
	assert.NoError(t, err)
}

func TestOsRemover_RemoveError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() { log.SetOutput(os.Stderr) }()

	r := &osRemover{}
	r.Remove("random_file")

	assert.Contains(t, buf.String(), "Could not remove random_file")
}

// removerMock is used for making sure Remove was called even if it was called asynchoronously
type removerMock struct {
	PathCalled string
	wg         sync.WaitGroup
}

func (r *removerMock) Remove(path string) {
	defer r.wg.Done()
	r.PathCalled = path
}

type pipelineMock struct {
	returnError string
	returnPath  string
	PathCalled  string
}

func (r *pipelineMock) Run(path string) (string, error) {
	if "" != r.returnError {
		return "", fmt.Errorf(r.returnError)
	}
	r.PathCalled = path

	return r.returnPath, nil
}

func TestRemoverPipeline_Run(t *testing.T) {
	mockRemover := &removerMock{}
	mockPipeline := &pipelineMock{returnPath: "expected_processed_path"}
	rp := RemoverPipeline{remover: mockRemover, pipeline: mockPipeline}

	mockRemover.wg.Add(1)
	processed, err := rp.Run("expected_original_path")
	mockRemover.wg.Wait()

	assert.NoError(t, err)
	assert.Equal(t, "expected_processed_path", processed)
	assert.Equal(t, "expected_original_path", mockPipeline.PathCalled)
	assert.Equal(t, "expected_original_path", mockRemover.PathCalled)
}

func TestRemoverPipeline_RunError(t *testing.T) {
	mockRemover := &removerMock{}
	mockPipeline := &pipelineMock{returnError: "expected_pipeline_error"}
	rp := RemoverPipeline{remover: mockRemover, pipeline: mockPipeline}

	processed, err := rp.Run("expected_original_path")

	assert.EqualError(t, err, "expected_pipeline_error")
	assert.Equal(t, "", mockRemover.PathCalled)
	assert.Equal(t, "", processed)
}
