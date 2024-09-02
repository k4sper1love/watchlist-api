package vcs

import (
	"fmt"
	"runtime/debug"
)

var version string

func init() {
	version = buildVersion()
}

func buildVersion() string {
	var (
		time     string
		revision string
		modified bool
	)

	bInfo, ok := debug.ReadBuildInfo()
	if !ok {
		version = "unknown"
	}

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

	if time == "" {
		time = "unknown-time"
	}

	if revision == "" {
		revision = "unknown-revision"
	}

	version = fmt.Sprintf("%s-%s", time, revision)
	if modified {
		version += "-dirty"
	}

	return version
}

func GetVersion() string {
	return version
}
