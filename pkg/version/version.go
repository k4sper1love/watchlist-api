package version

import "os"

func GetVersion() string {
	if version := os.Getenv("VERSION"); version != "" {
		return version
	} else {
		return "unknown"
	}
}
