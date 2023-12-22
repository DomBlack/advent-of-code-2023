package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// Context is the runner
type Context struct {
	context.Context                // base context
	log             zerolog.Logger // logger
	day             int            // day number
	part            int            // part number
	test            *testing.T     // If running as a test, this is the test
	saveOutput      bool           // record output to file
	isTest          bool           // is this run part of a test
}

func (c *Context) Ctx() context.Context {
	return c.Context
}

func (c *Context) Day() int {
	return c.day
}

// SaveOutput returns true if the output should be saved to a file
func (c *Context) SaveOutput() bool {
	if c == nil {
		return false
	}

	return c.saveOutput
}

// OutputFile returns the path to the output file for this day and part
func (c *Context) OutputFile(ext string) string {
	if c == nil {
		return ""
	}

	var dir string
	if c.isTest {
		dir = filepath.Join(repoDir, "internal", fmt.Sprintf("day%02d", c.day), "testdata")
	} else {
		dir = filepath.Join(repoDir, "outputs")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		c.log.Fatal().Err(err).Str("dir", dir).Msg("failed to create output directory")
	}

	if c.isTest {
		if c.test != nil {
			return filepath.Join(dir, strings.ToLower(fmt.Sprintf("part%02d_%s.%s", c.part, filepath.Base(c.test.Name()), ext)))
		} else {
			return filepath.Join(dir, fmt.Sprintf("part%02d.%s", c.part, ext))
		}
	} else {
		return filepath.Join(dir, fmt.Sprintf("day%02d_part%02d.%s", c.day, c.part, ext))
	}
}
