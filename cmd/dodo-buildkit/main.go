package main

import (
	"os"

	plugin "github.com/wabenet/dodo-buildkit"
)

func main() {
	os.Exit(plugin.RunMe())
}
