package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	uploadDir = "/var/data"
)

var authPassword string

func main() {
	// Get the auth password from the environment
	authPassword = os.Getenv("AUTH_PASSWORD")
	if authPassword == "" {
		fmt.Println("Warning: AUTH_PASSWORD not set in environment")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	http.HandleFunc("/upload", uploadHandler)

	port := ":8080"
	fmt.Printf("Server starting on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Printf("Error: Method not allowed - %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check Basic Auth
	_, password, ok := r.BasicAuth()
	if !ok || !validateCredentials(password) {
		fmt.Printf("Error: Unauthorized access attempt\n")
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		fmt.Printf("Error: Unable to parse form - %s\n", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("Error: Error retrieving file from form - %s\n", err)
		http.Error(w, "Error retrieving file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Printf("Error: Unable to create upload directory - %s\n", err)
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	// Create the destination file
	dst, err := os.Create(filepath.Join(uploadDir, header.Filename))
	if err != nil {
		fmt.Printf("Error: Unable to create the file for writing - %s\n", err)
		http.Error(w, "Unable to create the file for writing", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		fmt.Printf("Error: Unable to write file - %s\n", err)
		http.Error(w, "Unable to write file", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Success: File '%s' uploaded successfully\n", header.Filename)
	fmt.Fprintf(w, "File uploaded successfully")
}

func validateCredentials(password string) bool {
	return password == authPassword
}