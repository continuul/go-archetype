package version

import (
	"fmt"
	"strings"
)

var (
	// Do not manually edit! These will be filled in by the compiler.

	// GitCommit is the git commit of the release.
	GitCommit string
	// GitDescribe is the tags associated to the release.
	GitDescribe string

	// Version is the main version number that is being run at the moment.
	Version = "unknown"

	// VersionPrerelease is a pre-release marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = "unknown"
)

// GetHumanVersion composes the parts of the version in a way that's suitable
// for displaying to humans.
func GetHumanVersion() string {
	version := Version
	if GitDescribe != "" {
		version = GitDescribe
	}

	release := VersionPrerelease
	if GitDescribe == "" && release == "" {
		release = "dev"
	}

	if release != "" {
		if !strings.HasSuffix(version, "-"+release) {
			// if we tagged a prerelease version then the release is in the version already
			version += fmt.Sprintf("-%s", release)
		}
		if GitCommit != "" {
			version += fmt.Sprintf(" (%s)", GitCommit)
		}
	}

	// Strip off any single quotes added by the git information.
	return strings.Replace(version, "'", "", -1)
}
