package builder

import (
	"context"
	"fmt"

	cli "github.com/docker/cli/cli/command"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types/registry"
	docker "github.com/docker/docker/client"
	"github.com/wabenet/dodo-buildkit/internal/image"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
)

const name = "buildkit"

var _ builder.ImageBuilder = &Builder{}

type Builder struct {
	client docker.APIClient
}

func New() *Builder {
	return &Builder{}
}

func NewFromClient(client docker.APIClient) *Builder {
	return &Builder{client: client}
}

func (p *Builder) Type() plugin.Type {
	return builder.Type
}

func (p *Builder) Metadata() plugin.Metadata {
	return plugin.NewMetadata(builder.Type, name)
}

func (p *Builder) Init() (plugin.Config, error) {
	client, err := p.ensureClient()
	if err != nil {
		return nil, err
	}

	ping, err := client.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not reach docker host: %w", err)
	}

	return map[string]string{
		"client_version":  client.ClientVersion(),
		"host":            client.DaemonHost(),
		"api_version":     ping.APIVersion,
		"builder_version": fmt.Sprintf("%v", ping.BuilderVersion),
		"os_type":         ping.OSType,
		"experimental":    fmt.Sprintf("%t", ping.Experimental),
	}, nil
}

func (*Builder) Cleanup() {}

func (p *Builder) ensureClient() (docker.APIClient, error) {
	if p.client == nil {
		dockerCLI, err := cli.NewDockerCli(cli.WithBaseContext(context.Background()))
		if err != nil {
			return nil, fmt.Errorf("could not get docker config: %w", err)
		}

		if err := dockerCLI.Initialize(&cliflags.ClientOptions{}); err != nil {
			return nil, fmt.Errorf("could not get docker config: %w", err)
		}

		p.client = dockerCLI.Client()
	}

	return p.client, nil
}

func (p *Builder) CreateImage(config builder.BuildConfig, stream *plugin.StreamConfig) (string, error) {
	c, err := p.ensureClient()
	if err != nil {
		return "", err
	}

	// TODO: Don't do this twice
	dockerCLI, err := cli.NewDockerCli(cli.WithBaseContext(context.Background()))
	if err != nil {
		return "", fmt.Errorf("could not get docker config: %w", err)
	}

	creds, _ := dockerCLI.ConfigFile().GetAllCredentials()
	authConfigs := make(map[string]registry.AuthConfig, len(creds))
	for k, auth := range creds {
		authConfigs[k] = registry.AuthConfig(auth)
	}

	img, err := image.NewImage(c, authConfigs, config, stream)
	if err != nil {
		return "", fmt.Errorf("could not initialize builder client: %w", err)
	}

	imageID, err := img.Get()
	if err != nil {
		return "", fmt.Errorf("could not resolve image: %w", err)
	}

	return imageID, nil
}
