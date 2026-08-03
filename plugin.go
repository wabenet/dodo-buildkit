package plugin

import (
	moby "github.com/moby/moby/client"
	impl "github.com/wabenet/dodo-buildkit/internal/plugin/builder"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
)

func RunMe() int {
	m := plugin.Init()
	m.ServePlugins(NewImageBuilder())

	return 0
}

func IncludeMe(m plugin.Manager) {
	m.IncludePlugins(NewImageBuilder())
}

func NewImageBuilder() builder.ImageBuilder {
	return impl.New()
}

func NewImageBuilderWithDockerClient(c *moby.Client) *impl.Builder {
	return impl.NewFromClient(c)
}
