FROM golang:1.24-alpine

WORKDIR /app

# ビルドに必要なパッケージのインストール
RUN apk add --no-cache git

# 依存関係のコピーとダウンロード
COPY go.mod .

# 依存関係の初期化とダウンロード
RUN go mod download
RUN go mod tidy

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN go build -o main .

# アプリケーションの実行
CMD ["./main"]