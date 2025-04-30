package image

import (
	"io"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
	"golang.org/x/net/context"
)

const (
	ErrNoClient       ImageError = "client may not be nil"
	ErrMissingImageID ImageError = "build complete, but the server did not send an image id"
)

type ImageError string

func (e ImageError) Error() string {
	return string(e)
}

type Image struct {
	config      builder.BuildConfig
	client      Client
	authConfigs map[string]registry.AuthConfig
	session     session
	stream      *plugin.StreamConfig
}

type Client interface {
	DialHijack(context.Context, string, string, map[string][]string) (net.Conn, error)
	ImageList(context.Context, image.ListOptions) ([]image.Summary, error)
	ImageBuild(context.Context, io.Reader, types.ImageBuildOptions) (types.ImageBuildResponse, error)
}

func NewImage(
	client Client,
	authConfigs map[string]registry.AuthConfig,
	config builder.BuildConfig,
	stream *plugin.StreamConfig,
) (*Image, error) {
	if client == nil {
		return nil, ErrNoClient
	}

	session, err := prepareSession(config.Context)
	if err != nil {
		return nil, err
	}

	return &Image{
		client:      client,
		authConfigs: authConfigs,
		config:      config,
		session:     session,
		stream:      stream,
	}, nil
}
