package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"decentralstore/blockchain-service/internal/api"
	"decentralstore/blockchain-service/internal/infrastructure"
	"decentralstore/blockchain-service/internal/usecase"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	log.Println("Starting Blockchain Service...")

	// 環境変数の取得
	ethereumRPCURL := os.Getenv("ETHEREUM_RPC_URL")
	if ethereumRPCURL == "" {
		log.Println("ETHEREUM_RPC_URL not set, using default: http://localhost:8545")
		ethereumRPCURL = "http://localhost:8545"
	}

	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	if contractAddress == "" {
		log.Fatal("CONTRACT_ADDRESS environment variable is not set")
	}

	// Ethereumクライアントの初期化
	ethereumClient, err := infrastructure.NewEthereumClient(ethereumRPCURL)
	if err != nil {
		log.Fatalf("Failed to create Ethereum client: %v", err)
	}
	defer ethereumClient.Close()

	// スマートコントラクトの初期化
	contract, err := infrastructure.NewFileMetadataContract(common.HexToAddress(contractAddress), ethereumClient)
	if err != nil {
		log.Fatalf("Failed to create smart contract instance: %v", err)
	}

	// ユースケースの初期化
	blockchainService := usecase.NewBlockchainService(contract)

	// ハンドラーの初期化
	handler := api.NewBlockchainHandler(blockchainService)

	// ルーターの設定
	mux := http.NewServeMux()
	mux.HandleFunc("/store", handler.StoreMetadata)
	mux.HandleFunc("/metadata", handler.GetMetadata)
	mux.HandleFunc("/update", handler.UpdateMetadata)

	// HTTPサーバーの設定
	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	// サーバーを非同期で起動
	go func() {
		log.Printf("Server is listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// シグナル処理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// グレースフルシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
