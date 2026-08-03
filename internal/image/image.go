package image

import (
	"context"
	"io"
	"net"

	"github.com/moby/moby/api/types/registry"
	moby "github.com/moby/moby/client"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
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
	ImageList(context.Context, moby.ImageListOptions) (moby.ImageListResult, error)
	ImageBuild(context.Context, io.Reader, moby.ImageBuildOptions) (moby.ImageBuildResult, error)
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
