# Dual Server Example

一つの Go アプリケーションで gRPC と HTTP サーバーを同居させる最小構成の例。

## 構成

```
dual-server-example/
├── proto/           # Protocol Buffers定義
├── server/          # gRPC・HTTPサーバー実装
├── main.go          # メインアプリケーション
├── Makefile         # ビルド・実行スクリプト
└── go.mod           # Goモジュール定義
```

## 前提条件

- Go 1.24+
- Protocol Buffers Compiler (protoc)
- protoc-gen-go
- protoc-gen-go-grpc

## セットアップ

```bash
# 依存関係をインストール
go mod download

# Protocol Buffersツールをインストール
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 使用方法

### 開発実行

```bash
make dev
```

### ビルド・実行

```bash
make build
make run
```

### 手動実行

```bash
# Protocol Buffersのコード生成
make proto

# アプリケーション実行
go run main.go
```

## エンドポイント

### gRPC (localhost:9000)

```bash
# 書籍一覧取得
grpcurl -plaintext localhost:9000 books.BookService/ListBooks

# 特定の書籍取得
grpcurl -plaintext localhost:9000 books.BookService/GetBook -d '{"book_id": "book_1"}'

# 書籍作成
grpcurl -plaintext localhost:9000 books.BookService/CreateBook -d '{"book": {"title": "Go言語", "author": "山田太郎"}}'
```

### HTTP (localhost:8080)

```bash
# 書籍一覧取得
curl http://localhost:8080/books

# 特定の書籍取得
curl http://localhost:8080/books/book_1

# 書籍作成
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title": "Go言語", "author": "山田太郎"}'
```

## 特徴

- ✅ 最小構成
- ✅ grapi 非依存
- ✅ 一つのプロセスで両サーバー起動
- ✅ シンプルな実装
- ✅ 学習用途に最適

## 注意点

この実装は**学習・プロトタイプ用途**です。本番環境では以下を考慮してください：

- エラーハンドリングの強化
- セキュリティ対策
- パフォーマンス最適化
- 監視・ログ機能
- テストの追加
