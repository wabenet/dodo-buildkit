module github.com/wabenet/dodo-buildkit

go 1.26.0

require (
	github.com/docker/cli v29.5.2+incompatible
	github.com/docker/docker v28.5.2+incompatible
	github.com/hashicorp/go-hclog v1.6.3
	github.com/moby/buildkit v0.30.0
	github.com/morikuni/aec v1.1.0
	github.com/opencontainers/go-digest v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.11.1
	github.com/tonistiigi/fsutil v0.0.0-20251211185533-a2aa163d723f
	github.com/tonistiigi/units v0.0.0-20180711220420-6950e57a87ea
	github.com/tonistiigi/vt100 v0.0.0-20240514184818-90bafcd6abab
	github.com/wabenet/dodo-core v0.11.2
	golang.org/x/net v0.55.0
	golang.org/x/sync v0.20.0
	golang.org/x/time v0.15.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/containerd/v2 v2.2.3 // indirect
	github.com/containerd/continuity v0.4.5 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/typeurl/v2 v2.2.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.9.5 // indirect
	github.com/docker/go-connections v0.7.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fvbommel/sortorder v1.1.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0 // indirect
	github.com/hashicorp/go-plugin v1.8.0 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/moby/api v1.54.2 // indirect
	github.com/moby/moby/client v0.4.1 // indirect
	github.com/moby/patternmatcher v0.6.1 // indirect
	github.com/moby/sys/atomicwriter v0.1.0 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.68.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.68.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.68.0 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.opentelemetry.io/proto/otlp v1.10.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.51.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260406210006-6f92a3bedf2d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260406210006-6f92a3bedf2d // indirect
	google.golang.org/grpc v1.81.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
