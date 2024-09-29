package rest

import (
	"github.com/k4sper1love/watchlist-api/api"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// swaggerHandler handles requests to Swagger UI
func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	setSwaggerInfo(r)

	httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	).ServeHTTP(w, r)
}

// setSwaggerInfo configures the Host and Schemes for SwaggerInfo
func setSwaggerInfo(r *http.Request) {
	api.SwaggerInfo.Host = r.Host
	api.SwaggerInfo.Schemes = []string{"http", "https"}
}
