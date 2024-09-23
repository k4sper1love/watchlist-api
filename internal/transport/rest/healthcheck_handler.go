package rest

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"net/http"
	"time"
)

type SystemInfo struct {
	Environment string `json:"environment" example:"prod"`
	Uptime      string `json:"uptime" example:"3h 26m 30s"`
	LastChecked string `json:"last_checked" example:"2024-09-24T00:41:20+05:00"`
}

type Dependency struct {
	Status       string `json:"status" example:"up"`
	ResponseTime string `json:"response_time,omitempty" example:"48ms"`
}

type HealthcheckResponse struct {
	Status       string                `json:"status" example:"operational"`
	SystemInfo   SystemInfo            `json:"systemInfo"`
	Dependencies map[string]Dependency `json:"dependencies"`
}

// HealthCheck godoc
// @Summary Check API status
// @Description Check the API status. Returns status and system information.
// @Tags monitoring
// @Accept json
// @Produce json
// @Success 200 {object} HealthcheckResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Router /healthcheck [get]
func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	systemInfo := SystemInfo{
		Environment: config.Env,
		Uptime:      getUptime(startServerTime),
		LastChecked: time.Now().Format(time.RFC3339),
	}

	dependencies := checkDependencies()

	response := HealthcheckResponse{
		Status:       getAPIStatus(dependencies),
		SystemInfo:   systemInfo,
		Dependencies: dependencies,
	}

	writeJSON(w, r, http.StatusOK, envelope{"healthcheck": response})
}

func checkDependencies() map[string]Dependency {
	dependencies := make(map[string]Dependency)

	// Check database health
	dbStatus, dbResponseTime := checkDB()
	dependencies["database"] = Dependency{
		Status:       dbStatus,
		ResponseTime: dbResponseTime,
	}

	return dependencies
}

func checkDB() (string, string) {
	start := time.Now()

	if err := postgres.PingDB(); err != nil {
		return "down", ""
	}

	responseTime := time.Since(start).Milliseconds()

	return "up", fmt.Sprintf("%dms", responseTime)
}

func getAPIStatus(dependencies map[string]Dependency) string {
	for _, dependency := range dependencies {
		if dependency.Status != "up" {
			return "degraded"
		}
	}
	return "operational"
}

func getUptime(startTime time.Time) string {
	duration := time.Since(startTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}
