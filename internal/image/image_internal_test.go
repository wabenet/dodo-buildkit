package image

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/stretchr/testify/assert"
	api "github.com/wabenet/dodo-core/api/build/v1alpha2"
	"golang.org/x/net/context"
)

func fakeImage(t *testing.T, config *api.BuildConfig) *Image {
	return &Image{
		client:  &fakeImageClient{t: t, willBuildAs: "NewImageID"},
		config:  config,
		session: &fakeSession{},
	}
}

type fakeImageClient struct {
	t           *testing.T
	willBuildAs string
}

func (client *fakeImageClient) DialHijack(
	_ context.Context, _ string, _ string, _ map[string][]string,
) (net.Conn, error) {
	return nil, nil
}

func (client *fakeImageClient) ImageList(
	_ context.Context, _ image.ListOptions,
) ([]image.Summary, error) {
	return []image.Summary{}, nil
}

func (client *fakeImageClient) BuildCancel(_ context.Context, _ string) error {
	return nil
}

func (client *fakeImageClient) ImageBuild(
	_ context.Context, _ io.Reader, _ types.ImageBuildOptions,
) (types.ImageBuildResponse, error) {
	buildResult := types.BuildResult{ID: client.willBuildAs}
	auxJSON, err := json.Marshal(buildResult)
	assert.Nil(client.t, err)

	rawJSON := json.RawMessage(auxJSON)
	message := jsonmessage.JSONMessage{
		ID:     "moby.image.id",
		Stream: "hello world",
		Aux:    &rawJSON,
	}
	response, err := json.Marshal(message)
	assert.Nil(client.t, err)

	return types.ImageBuildResponse{Body: io.NopCloser(bytes.NewReader(response))}, nil
}
