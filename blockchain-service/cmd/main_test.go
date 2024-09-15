package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// テスト全体の前処理
	setup()
	// テストの実行
	code := m.Run()
	// テスト全体の後処理
	teardown()
	// テスト結果でOSを終了
	os.Exit(code)
}

func setup() {
	// テストに必要な環境変数を設定
	os.Setenv("ETHEREUM_RPC_URL", "http://localhost:8545")
	os.Setenv("CONTRACT_ADDRESS", "0x1234567890123456789012345678901234567890")
}

func teardown() {
	// 環境変数をクリア
	os.Unsetenv("ETHEREUM_RPC_URL")
	os.Unsetenv("CONTRACT_ADDRESS")
}

func TestServerStart(t *testing.T) {
	// mainを非同期で実行
	go main()

	// サーバーの起動を待つ
	time.Sleep(100 * time.Millisecond)

	// サーバーが起動しているか確認
	resp, err := http.Get("http://localhost:8082/metadata?fileID=test")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status OK or Internal Server Error, got %v", resp.Status)
	}
}

func TestMissingContractAddress(t *testing.T) {
	// CONTRACT_ADDRESSを一時的に削除
	origContractAddress := os.Getenv("CONTRACT_ADDRESS")
	os.Unsetenv("CONTRACT_ADDRESS")
	defer os.Setenv("CONTRACT_ADDRESS", origContractAddress)

	// ログ出力をキャプチャ
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// mainを実行
	main()

	// ログ出力を確認
	logOutput := buf.String()
	if !strings.Contains(logOutput, "CONTRACT_ADDRESS environment variable is not set") {
		t.Errorf("Expected log message about missing CONTRACT_ADDRESS, got: %s", logOutput)
	}
}
