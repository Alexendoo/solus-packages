// +build linux darwin

package packages

import (
	"os"
)

func CachePath() string {
	xdg := os.Getenv("XDG_CACHE_HOME")
	if xdg != "" {
		return xdg
	}

	return os.ExpandEnv("$HOME/.cache")
}
