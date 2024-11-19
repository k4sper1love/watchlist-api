package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// uploadImageHandler handles image upload.
func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		badRequestResponse(w, r, fmt.Errorf("error parsing the form: %v", err))
		return
	}

	imageFile, _, err := r.FormFile("image")
	if err != nil {
		badRequestResponse(w, r, fmt.Errorf("error receiving the file: %v", err))
		return
	}
	defer imageFile.Close()

	filename := fmt.Sprintf("%d_%s.jpg", time.Now().UnixNano(), generateString(5))
	filePath := filepath.Join("static/images", filename)

	out, err := os.Create(filePath)
	if err != nil {
		serverErrorResponse(w, r, fmt.Errorf("error creating the file: %v", err))
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, imageFile); err != nil {
		serverErrorResponse(w, r, fmt.Errorf("error copying file: %v", err))
		return
	}

	imageURL := fmt.Sprintf("http://%s/images/%s", r.Host, filename)

	writeJSON(w, r, http.StatusCreated, envelope{"image_url": imageURL})
}

// getImageHandler handles image retrieval requests by file name.
func getImageHandler(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]

	pathToFile := filepath.Join("static/images", filename)

	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		badRequestResponse(w, r, fmt.Errorf("file not found"))
		return
	}

	http.ServeFile(w, r, pathToFile)
}

func setDefaultImage(r *http.Request, f *models.Film) {
	if f.ImageURL == "" {
		f.ImageURL = fmt.Sprintf("http://%s/images/default.png", r.Host)
	}
}
