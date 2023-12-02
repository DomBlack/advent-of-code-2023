package runner

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var repoDir string

func init() {
	// Setup the base logging
	log.Logger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Caller().Logger()

	// Find the repo root
	var err error
	repoDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	i := 0
	for {
		if _, err := os.Stat(repoDir + "/go.mod"); err == nil {
			break
		}

		repoDir = repoDir + "/.."
		i++
		if i > 10 {
			panic("failed to find repo root")
		}
	}

	repoDir = filepath.Clean(repoDir)
}
