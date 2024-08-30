package vcs

import (
	"fmt"
	"log"
	"runtime/debug"
)

func Version() string {
	var (
		time     string
		revision string
		modified bool
	)

	bInfo, ok := debug.ReadBuildInfo()
	if !ok {
		log.Println("Build info not available")
		return ""
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

	if modified {
		return fmt.Sprintf("%s-%s-dirty", time, revision)
	}

	return fmt.Sprintf("%s-%s", time, revision)
}
