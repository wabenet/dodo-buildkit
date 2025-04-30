package image

import (
	"testing"

	controlapi "github.com/moby/buildkit/api/services/control"
	"github.com/stretchr/testify/assert"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
)

func TestBuildImage(t *testing.T) {
	displayCh := make(chan *controlapi.StatusResponse)
	defer close(displayCh)

	image := fakeImage(t, builder.BuildConfig{
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

	image := fakeImage(t, builder.BuildConfig{
		InlineDockerfile: []string{"FROM scratch"},
	})
	result, err := image.runBuild(&contextData{
		remote:         "client-session",
		dockerfileName: "Dockerfile",
	}, displayCh)
	assert.Nil(t, err)
	assert.Equal(t, "NewImageID", result)
}
