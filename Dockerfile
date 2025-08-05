FROM golang:1.24-alpine AS builder

WORKDIR /app

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# Protocol Buffersツールをインストール
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# ソースコードをコピー
COPY . .

# Protocol Buffersのコード生成
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/books.proto

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/dual-server main.go

# 実行用イメージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドしたアプリケーションをコピー
COPY --from=builder /app/dual-server .

# ポートを公開
EXPOSE 9000 8080

# アプリケーションを起動
CMD ["./dual-server"]
