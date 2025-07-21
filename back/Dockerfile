FROM golang:1.20.2 AS builder

# 作業ディレクトリを作成
RUN mkdir /app
WORKDIR /app

# Goモジュールを有効化
ENV GO111MODULE=on

COPY . .

# 依存関係を整理し、必要なものだけを残す
RUN go mod tidy

RUN go build -o main .

FROM alpine as dev
WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE 8080

# エントリーポイントを設定
CMD ["./main"]