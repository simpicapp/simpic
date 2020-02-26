package internal

import (
	"fmt"
	"strings"
)

var (
	// Output of `git describe --tags --always`
	GitTag string

	// Output of `git rev-parse --short HEAD`
	GitSHA string
)

func GetVersionString() string {
	if GitTag == "" {
		return "Simpic"
	}

	if GitTag == GitSHA {
		return fmt.Sprintf("Simpic/pre-release (%s)", GitSHA)
	}

	parts := strings.SplitN(GitTag, "-", 2)
	if len(parts) > 1 {
		// Some kind of untagged version
		return fmt.Sprintf("Simpic/%s+dev (%s)", parts[0], GitSHA)
	}

	// Actual tagged commit
	return fmt.Sprintf("Simpic/%s (%s)", parts[0], GitSHA)
}
