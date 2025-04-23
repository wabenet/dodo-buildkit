package image

import (
	"testing"

	controlapi "github.com/moby/buildkit/api/services/control"
	"github.com/stretchr/testify/assert"
	api "github.com/wabenet/dodo-core/api/build/v1alpha2"
)

func TestBuildImage(t *testing.T) {
	displayCh := make(chan *controlapi.StatusResponse)
	defer close(displayCh)

	image := fakeImage(t, &api.BuildConfig{
		Context: "./test",
	})
	result, err := image.runBuild(&contextData{
		remote:         "client-session",
		dockerfileName: "Dockerfile",
	}, displayCh)
	assert.Nil(t, err)
	assert.Equal(t, "NewImageID", result)
}

func TestBuildInlineImage(t *testing.T) {
	displayCh := make(chan *controlapi.StatusResponse)
	defer close(displayCh)

	image := fakeImage(t, &api.BuildConfig{
		InlineDockerfile: []string{"FROM scratch"},
	})
	result, err := image.runBuild(&contextData{
		remote:         "client-session",
		dockerfileName: "Dockerfile",
	}, displayCh)
	assert.Nil(t, err)
	assert.Equal(t, "NewImageID", result)
}
