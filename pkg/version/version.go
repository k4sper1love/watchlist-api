package version

import "os"

func GetVersion() string {
	return os.Getenv("VERSION")
}
