# 階段一：編譯階段
FROM golang:1.24-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 複製 Go 套件清單並下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製所有原始碼
COPY . .

# 編譯成二進位執行檔（針對 Linux 環境優化）
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service main.go

# 階段二：運行階段（使用極小的 alpine 系統）
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 從編譯階段把執行檔 copy 過來
COPY --from=builder /app/user-service .

# 暴露 Cloud Run 預設的 8080 Port
EXPOSE 8080

# 啟動服務
CMD ["./user-service"]