package image

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/containerd/containerd/v2/pkg/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	controlapi "github.com/moby/buildkit/api/services/control"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/moby/moby/api/types/build"
	"github.com/moby/moby/api/types/jsonstream"
	moby "github.com/moby/moby/client"
	"github.com/wabenet/dodo-buildkit/internal/progress"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func (img *Image) Get() (string, error) {
	if img.config.ForceRebuild || len(img.config.ImageName) == 0 {
		return img.Build()
	}

	imgs, err := img.client.ImageList(
		context.Background(),
		moby.ImageListOptions{
			Filters: make(moby.Filters).Add("reference", img.config.ImageName),
		},
	)
	if err != nil || len(imgs.Items) == 0 {
		return img.Build()
	}

	return imgs.Items[0].ID, nil
}

func (image *Image) Build() (string, error) {
	contextData, err := prepareContext(image.config, image.session)
	if err != nil {
		return "", err
	}

	defer contextData.cleanup()

	imageID := ""
	displayCh := make(chan *controlapi.StatusResponse)

	eg, _ := errgroup.WithContext(appcontext.Context())

	eg.Go(func() error {
		return image.session.Run(
			context.TODO(),
			func(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error) {
				return image.client.DialHijack(ctx, "/session", proto, meta)
			},
		)
	})

	if image.stream != nil {
		eg.Go(func() error {
			ctx := context.TODO()
			t := progress.NewPrinter(
				image.stream.Stderr,
				int(image.stream.TerminalHeight),
				int(image.stream.TerminalWidth),
			)

			tickerTimeout := 150 * time.Millisecond
			displayTimeout := 100 * time.Millisecond

			if v := os.Getenv("TTY_DISPLAY_RATE"); v != "" {
				if r, err := strconv.ParseInt(v, 10, 64); err == nil {
					tickerTimeout = time.Duration(r) * time.Millisecond
					displayTimeout = time.Duration(r) * time.Millisecond
				}
			}

			var done bool
			ticker := time.NewTicker(tickerTimeout)
			defer ticker.Stop()

			displayLimiter := rate.NewLimiter(rate.Every(displayTimeout), 1)

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
				case ss, ok := <-displayCh:
					if ok {
						t.Update(ss)
					} else {
						done = true
					}
				}

				if done {
					t.Print(true)
					t.PrintErrorLogs()

					return nil
				}

				if displayLimiter.Allow() {
					ticker.Stop()
					ticker = time.NewTicker(tickerTimeout)
					t.Print(false)
				}
			}
		})
	}

	eg.Go(func() error {
		defer func() {
			close(displayCh)
			image.session.Close()
		}()

		imageID, err = image.runBuild(contextData, displayCh)

		return err
	})

	err = eg.Wait()
	if err != nil {
		return "", fmt.Errorf("error during build: %w", err)
	}

	if imageID == "" {
		return "", ErrMissingImageID
	}

	return imageID, nil
}

func (image *Image) runBuild(contextData *contextData, displayCh chan *controlapi.StatusResponse) (string, error) {
	args := map[string]*string{}
	for _, arg := range image.config.Arguments {
		args[arg.Key] = &arg.Value
	}

	var tags []string
	if image.config.ImageName != "" {
		tags = append(tags, image.config.ImageName)
	}

	response, err := image.client.ImageBuild(
		context.Background(),
		nil,
		moby.ImageBuildOptions{
			Tags:           tags,
			SuppressOutput: false,
			NoCache:        image.config.NoCache,
			Remove:         true,
			ForceRemove:    true,
			PullParent:     image.config.ForcePull,
			Dockerfile:     contextData.dockerfileName,
			BuildArgs:      args,
			AuthConfigs:    image.authConfigs,
			Version:        build.BuilderBuildKit,
			RemoteContext:  contextData.remote,
			SessionID:      image.session.ID(),
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not build image: %w", err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	var imageID string

	for {
		var msg jsonstream.Message
		if err := decoder.Decode(&msg); err != nil {
			if errors.Is(err, io.EOF) {
				return imageID, nil
			}

			return "", fmt.Errorf("could not decode JSON message: %w", err)
		}

		if msg.Error != nil {
			return "", msg.Error
		}

		if msg.Aux == nil {
			continue
		}

		switch msg.ID {
		case "moby.image.id":
			var result build.Result

			err := json.Unmarshal(*msg.Aux, &result)
			if err != nil {
				continue
			}

			log.L().Debug("received image id", "id", result.ID)

			imageID = result.ID

		case "moby.buildkit.trace":
			if image.stream == nil {
				continue
			}

			var protoResponse []byte

			err := json.Unmarshal(*msg.Aux, &protoResponse)
			if err != nil {
				log.L().Debug("unmarshal error", "aux", msg.Aux, "error", err)

				continue
			}

			var resp controlapi.StatusResponse

			err = proto.Unmarshal(protoResponse, &resp)
			if err != nil {
				log.L().Debug("unmarshal error", "aux", msg.Aux, "error", err)

				continue
			}

			log.L().Debug("received build trace", "status", &resp)

			displayCh <- &resp
		}
	}
}
