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
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}

	return fmt.Sprintf("%ds", seconds)
}
