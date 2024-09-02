/*
Package vcs provides functions to retrieve
and format version control system metadata for the application.

This package extracts VCS information like commit time, revision,
and modification status, and stores it in a version string accessible via GetVersion.
*/

package vcs

import (
	"fmt"
	"runtime/debug"
)

// version is a package-level variable for storing the version string.
var version string

// init initializes the version by calling the buildVersion function.
// This function is performed automatically when the package is imported.
func init() {
	buildVersion()
}

// buildVersion constructs a version string from VCS metadata
// such as commit time, revision, and modification status.
// If any of this information is unavailable,
// default values like "unknown-time" and "unknown-revision" are used instead.
func buildVersion() {
	var (
		time     string // Variable to store the commit time from the VCS.
		revision string // Variable to store the commit revision from the VCS.
		modified bool   // Flag indicating if the working directory was modified.
	)

	// Retrieve the build information, which includes VCS metadata.
	bInfo, ok := debug.ReadBuildInfo()
	// Check if the build information is available.
	// Return a default version string if unavailable.
	if !ok {
		version = "unknown-time-unknown-revision"
		return
	}

	// Loop through the build settings to extract VCS metadata.
	for _, bSetting := range bInfo.Settings {
		switch bSetting.Key {
		case "vcs.time":
			time = bSetting.Value
		case "vcs.revision":
			revision = bSetting.Value
		case "vcs.modified":
			if bSetting.Value == "true" {
				modified = true
			}
		}
	}

	// Check if the time variable is empty.
	// Assign a default value if no commit time is available.
	if time == "" {
		time = "unknown-time"
	}

	// Check if the revision variable is empty.
	// Assign a default value if no commit revision is available.
	if revision == "" {
		revision = "unknown-revision"
	}

	// Construct the version string with time and revision.
	version = fmt.Sprintf("%s-%s", time, revision)

	// Check if the modified flag is true.
	// Append "-dirty" to the version string if there were uncommitted changes.
	if modified {
		version += "-dirty"
	}
}

// GetVersion returns the version string that was constructed during package initialization.
// The version includes the commit time and revision from the VCS, and if applicable,
// a "-dirty" suffix if there were uncommitted changes at build time.
func GetVersion() string {
	return version
}
