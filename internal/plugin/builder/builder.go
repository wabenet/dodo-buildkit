package builder

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/context/docker"
	"github.com/docker/cli/cli/context/store"
	"github.com/moby/moby/api/types/registry"
	moby "github.com/moby/moby/client"
	"github.com/wabenet/dodo-buildkit/internal/image"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
)

const (
	name = "buildkit"

	defaultDockerContext = "default"
	defaultDockerHost    = "unix:///var/run/docker.sock"
)

var _ builder.ImageBuilder = &Builder{}

type Builder struct {
	client moby.APIClient
	config *configfile.ConfigFile
}

func New() *Builder {
	return &Builder{}
}

func NewFromClient(client moby.APIClient) *Builder {
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

	ping, err := client.Ping(context.Background(), moby.PingOptions{})
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

func (p *Builder) CreateImage(config builder.BuildConfig, stream *plugin.StreamConfig) (string, error) {
	client, err := p.ensureClient()
	if err != nil {
		return "", err
	}

	creds, _ := p.config.GetAllCredentials()
	authConfigs := make(map[string]registry.AuthConfig, len(creds))

	for k, auth := range creds {
		authConfigs[k] = registry.AuthConfig(auth)
	}

	img, err := image.NewImage(client, authConfigs, config, stream)
	if err != nil {
		return "", fmt.Errorf("could not initialize builder client: %w", err)
	}

	imageID, err := img.Get()
	if err != nil {
		return "", fmt.Errorf("could not resolve image: %w", err)
	}

	return imageID, nil
}

func (p *Builder) ensureClient() (moby.APIClient, error) { //nolint:ireturn
	if p.client == nil {
		p.config = config.LoadDefaultConfigFile(nil)

		endpoint, err := getDockerEndpoint(p.config)
		if err != nil {
			return nil, fmt.Errorf("could not get docker endpoint: %w", err)
		}

		opts, err := endpoint.ClientOpts()
		if err != nil {
			return nil, fmt.Errorf("could not get endpoint options: %w", err)
		}

		client, err := moby.New(opts...)
		if err != nil {
			return nil, fmt.Errorf("could not create moby client: %w", err)
		}

		p.client = client
	}

	return p.client, nil
}

func getDockerEndpoint(configFile *configfile.ConfigFile) (docker.Endpoint, error) {
	ctxName := getContextName(configFile)

	if ctxName == defaultDockerContext {
		return defaultDockerEndpoint()
	}

	return dockerEndpointFromContext(ctxName)
}

func getContextName(configFile *configfile.ConfigFile) string {
	if os.Getenv(moby.EnvOverrideHost) != "" {
		return defaultDockerContext
	}

	if ctxName := os.Getenv("DOCKER_CONTEXT"); ctxName != "" {
		return ctxName
	}

	if configFile.CurrentContext != "" {
		return configFile.CurrentContext
	}

	return defaultDockerContext
}

func defaultDockerEndpoint() (docker.Endpoint, error) {
	endpoint := docker.Endpoint{
		EndpointMeta: docker.EndpointMeta{
			Host:          defaultDockerHost,
			SkipTLSVerify: false,
		},
	}

	if override := os.Getenv(moby.EnvOverrideHost); override != "" {
		// The original Docker CLI uses a whole lot of logic to infer a valid endpoint from the env var,
		// including lots of default values and so on.
		// We are just assuming the user passes the endpoint in the correct format and hope for the best.
		endpoint.Host = override
	}

	return endpoint, nil
}

func dockerEndpointFromContext(name string) (docker.Endpoint, error) {
	ctxStore := store.New(
		config.ContextStoreDir(),
		store.NewConfig(
			func() any { return &dockerContext{} },
			[]store.NamedTypeGetter{
				store.EndpointTypeGetter(docker.DockerEndpoint, func() any { return &docker.EndpointMeta{} }),
			}...,
		),
	)

	ctxMeta, err := ctxStore.GetMetadata(name)
	if err != nil {
		return docker.Endpoint{}, fmt.Errorf("could not get context metadata: %w", err)
	}

	epMeta, err := docker.EndpointFromContext(ctxMeta)
	if err != nil {
		return docker.Endpoint{}, fmt.Errorf("could not get endpoint from context: %w", err)
	}

	endpoint, err := docker.WithTLSData(ctxStore, name, epMeta)
	if err != nil {
		return docker.Endpoint{}, fmt.Errorf("could not create endpoint: %w", err)
	}

	return endpoint, nil
}

type dockerContext struct {
	Description      string
	AdditionalFields map[string]any
}
