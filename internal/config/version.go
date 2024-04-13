package config

import (
	"fmt"
	"runtime/debug"
)

// Helper to read version from VCS
func version() string {
	var (
		time     string
		modified bool
		revision string
	)

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}

			case "vcs.revision":
				revision = s.Value

			case "vcs.time":
				time = s.Value
			}
		}
	}

	if modified {
		return fmt.Sprintf("%s-%s-dirty", time, revision)
	}

	return revision
}
