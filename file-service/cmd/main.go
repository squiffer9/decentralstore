package main

import (
	"log"
	"net/http"
	"os"

	"decentralstore/file-service/internal/api"
	"decentralstore/file-service/internal/infrastructure"
	"decentralstore/file-service/internal/usecase"
)

func main() {
	log.Println("Starting File Service...")

	ipfsAPI := os.Getenv("IPFS_API_URL")
	if ipfsAPI == "" {
		ipfsAPI = "localhost:5001"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	storageClient, err := infrastructure.NewStorageClient(ipfsAPI, redisURL)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}

	router := SetupRoutes(storageClient)

	log.Fatal(http.ListenAndServe(":8081", router))
}

func SetupRoutes(storageClient *infrastructure.StorageClient) http.Handler {
	fileHandler := CreateFileHandler(storageClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", fileHandler.UploadFile)
	mux.HandleFunc("/download", fileHandler.DownloadFile)
	mux.HandleFunc("/delete", fileHandler.DeleteFile)

	return mux
}

func CreateFileHandler(storageClient *infrastructure.StorageClient) *api.FileHandler {
	fileUseCase := usecase.NewFileUseCase(storageClient)
	return api.NewFileHandler(fileUseCase)
}
