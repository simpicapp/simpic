package internal

import (
	"fmt"
)

const (
	Version = "1.0+dev"
)

var (
	// Output of `git rev-parse --short HEAD`
	GitSHA string
)

func GetVersionString() string {
	if GitSHA == "" {
		return fmt.Sprintf("Simpic/%s", Version)
	}

	return fmt.Sprintf("Simpic/%s (%s)", Version, GitSHA)
}
