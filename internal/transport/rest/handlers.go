package rest

import (
	"fmt"
	"log"
	"net/http"
)

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("healthcheckHandler serving:", r.URL.Path, r.Host)
	message := fmt.Sprintf("Wishlist API is working! URL Path: %s, Host: %s", r.URL.Path, r.Host)
	writeJSON(w, r, http.StatusOK, envelope{"message": message})
}
