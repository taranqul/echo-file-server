package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cant read body of request!", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Echoed: %s", string(body))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "I'm healthy!")

}

func echoFileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Can't parse file: "+err.Error(), http.StatusBadRequest)
	}

	saveDir := "uploads"
	savePath := filepath.Join(saveDir, header.Filename)

	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Can't create file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	tee := io.TeeReader(file, dst)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", header.Filename))

	if _, err = io.Copy(w, tee); err != nil {
		http.Error(w, "Can't echo file: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /echo", echoHandler)
	mux.HandleFunc("POST /echo-file", echoFileHandler)
	mux.HandleFunc("GET /health", healthHandler)
	saveDir := "uploads"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		log.Fatal("Can't create dir: " + saveDir)
		return
	}
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
