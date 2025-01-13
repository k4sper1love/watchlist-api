package postgres

import (
	"fmt"
	"strconv"
	"strings"
)

func parseRangeOrExact(value string) (float64, float64, int, error) {
	if value == "" {
		return 0, 0, -1, nil
	}
	if strings.Contains(value, "-") {
		parts := strings.Split(value, "-")
		if len(parts) != 2 {
			return 0, 0, 0, fmt.Errorf("invalid range format")
		}
		minValue, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return 0, 0, 0, err
		}
		maxValue, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return 0, 0, 0, err
		}
		return minValue, maxValue, 1, nil
	}

	exact, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	return exact, exact, 0, nil
}

func parseYearOrRange(value string) (int, int, int, error) {
	if value == "" {
		return 0, 0, -1, nil
	}
	if strings.Contains(value, "-") {
		parts := strings.Split(value, "-")
		if len(parts) != 2 {
			return 0, 0, 0, fmt.Errorf("invalid year range format")
		}
		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return 0, 0, 0, err
		}
		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return 0, 0, 0, err
		}
		return start, end, 1, nil
	}

	year, err := strconv.Atoi(value)
	if err != nil {
		return 0, 0, 0, err
	}
	return year, year, 0, nil
}
