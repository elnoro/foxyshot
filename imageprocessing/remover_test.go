package imageprocessing

import (
	"fmt"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsRemover_Remove(t *testing.T) {
	r := &osRemover{}

	f, err := ioutil.TempFile("testdata/", "tempfile")
	if err != nil {
		assert.FailNow(t, "Cannot run this test - no file was created")
	}
	defer f.Close()

	err = r.Remove(f.Name())

	assert.NoError(t, err)
	assert.False(t, assert.FileExists(new(testing.T), f.Name()), "Temp file not removed by remover")

	err = r.Remove("random_file")
	assert.Error(t, err)
}

// removerMock is used for making sure Remove was called even if it was called asynchoronously
type removerMock struct {
	PathCalled string
	wg         sync.WaitGroup
}

func (r *removerMock) Remove(path string) error {
	defer r.wg.Done()
	r.PathCalled = path

	return nil
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

var examples = []struct {
	p pipelineMock
	r removerMock
}{}

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
