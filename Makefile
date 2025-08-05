.PHONY: proto build run clean test

# Protocol Buffersのコード生成
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/books.proto

# アプリケーションをビルド
build: proto
	go build -o bin/dual-server main.go

# アプリケーションを実行
run: build
	./bin/dual-server

# 依存関係をダウンロード
deps:
	go mod download

# テスト実行
test:
	go test ./...

# クリーンアップ
clean:
	rm -rf bin/
	rm -f proto/*.pb.go

# 開発用（コード生成 + 実行）
dev: proto
	go run main.go
