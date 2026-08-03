package image

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"testing"

	"github.com/moby/moby/api/types/build"
	"github.com/moby/moby/api/types/jsonstream"
	moby "github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
)

func fakeImage(t *testing.T, config builder.BuildConfig) *Image {
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
	_ context.Context, _ moby.ImageListOptions,
) (moby.ImageListResult, error) {
	return moby.ImageListResult{}, nil
}

func (client *fakeImageClient) BuildCancel(_ context.Context, _ string) error {
	return nil
}

func (client *fakeImageClient) ImageBuild(
	_ context.Context, _ io.Reader, _ moby.ImageBuildOptions,
) (moby.ImageBuildResult, error) {
	buildResult := build.Result{ID: client.willBuildAs}
	auxJSON, err := json.Marshal(buildResult)
	assert.Nil(client.t, err)

	rawJSON := json.RawMessage(auxJSON)
	message := jsonstream.Message{
		ID:     "moby.image.id",
		Stream: "hello world",
		Aux:    &rawJSON,
	}
	response, err := json.Marshal(message)
	assert.Nil(client.t, err)

	return moby.ImageBuildResult{Body: io.NopCloser(bytes.NewReader(response))}, nil
}
