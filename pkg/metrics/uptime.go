package metrics

import (
	"fmt"
	"time"
)

var startTime time.Time

func InitUptime() {
	startTime = time.Now()
}

func GetUptime() float64 {
	return time.Since(startTime).Seconds()
}

func GetUptimeFormat() string {
	duration := time.Since(startTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}
