package config

import (
	"os"
)

func CachePath() string {
	return os.ExpandEnv("$LOCALAPPDATA\\Solus Packages")
}
