package api

import (
	"encoding/json"
	"net/http"

	"decentralstore/blockchain-service/internal/domain"
	"decentralstore/blockchain-service/internal/usecase"
)

type BlockchainHandler struct {
	service usecase.BlockchainService
}

func NewBlockchainHandler(service usecase.BlockchainService) *BlockchainHandler {
	return &BlockchainHandler{service: service}
}

func (h *BlockchainHandler) StoreMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var metadata domain.FileMetadata
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.StoreMetadata(r.Context(), &metadata); err != nil {
		http.Error(w, "Failed to store metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Metadata stored successfully"})
}

func (h *BlockchainHandler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileID := r.URL.Query().Get("fileID")
	if fileID == "" {
		http.Error(w, "Missing fileID parameter", http.StatusBadRequest)
		return
	}

	metadata, err := h.service.GetMetadata(r.Context(), fileID)
	if err != nil {
		http.Error(w, "Failed to get metadata", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

func (h *BlockchainHandler) UpdateMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileID := r.URL.Query().Get("fileID")
	if fileID == "" {
		http.Error(w, "Missing fileID parameter", http.StatusBadRequest)
		return
	}

	var updateRequest struct {
		IsDeleted bool `json:"isDeleted"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateMetadata(r.Context(), fileID, updateRequest.IsDeleted); err != nil {
		http.Error(w, "Failed to update metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Metadata updated successfully"})
}
