package main

import (
	"fmt"
	"net/http"
	"os"
	"smartraveller-api/api"
)

func main() {
	http.HandleFunc("/advisory", api.GetAdvisory)
	http.HandleFunc("/advisories", api.GetAdvisories)

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Server started successfully")
}
