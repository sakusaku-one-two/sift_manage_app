# Multi-stage build for Go application
# Multi-stage: マルチステージ、build: ビルド、application: アプリケーション

# Build stage
# build: ビルド、stage: ステージ
FROM golang:1.21-alpine AS builder  # golang: Go言語、alpine: Alpine Linux、builder: ビルダー

# Set working directory
# set: 設定する、working: 作業、directory: ディレクトリ
WORKDIR /app

# Install build dependencies
# install: インストール、dependencies: 依存関係（複数形）
RUN apk add --no-cache git ca-certificates tzdata  # apk: Alpine Package Keeper、no-cache: キャッシュなし、git: バージョン管理、ca-certificates: CA証明書、tzdata: タイムゾーンデータ

# Copy go mod files
# copy: コピー、mod: module（モジュール）、files: ファイル（複数形）
COPY go.mod go.sum ./

# Download dependencies
# download: ダウンロード
RUN go mod download  # download: ダウンロード、依存関係をダウンロード

# Copy source code
# source: ソース、code: コード
COPY . .

# Build the application
# application: アプリケーション
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app  # cgo: C言語バインディング、enabled: 有効、goos: Go OS、linux: Linux OS、installsuffix: インストールサフィックス、output: 出力

# Production stage
# production: 本番、stage: ステージ
FROM alpine:latest AS production  # latest: 最新版、production: 本番

# Install runtime dependencies
# runtime: 実行時、dependencies: 依存関係
RUN apk --no-cache add ca-certificates curl  # curl: HTTP クライアント

# Create non-root user for security
# create: 作成する、non-root: 非ルート、user: ユーザー、security: セキュリティ
RUN addgroup -g 1001 -S appgroup && \  # addgroup: グループ追加、-g: group ID、-S: system group
    adduser -u 1001 -S appuser -G appgroup  # adduser: ユーザー追加、-u: user ID、-G: group

# Set working directory
WORKDIR /app

# Copy binary from builder stage
# binary: バイナリ、from: から、builder: ビルダー
COPY --from=builder /app/main .

# Copy configuration files if needed
# configuration: 設定、files: ファイル、needed: 必要な
# COPY --from=builder /app/config ./config

# Create logs directory
# logs: ログ（複数形）
RUN mkdir -p /app/logs && chown -R appuser:appgroup /app  # mkdir: ディレクトリ作成、-p: parent directories、chown: change owner（所有者変更）、-R: recursive（再帰的）

# Switch to non-root user
# switch: 切り替える、to: に
USER appuser  # user: ユーザー

# Expose application port
# expose: 公開する、port: ポート
EXPOSE 8080

# Health check endpoint
# health: 健康、check: 確認、endpoint: エンドポイント
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1  # cmd: コマンド、-f: fail silently、exit: 終了

# Set environment variables
# environment: 環境、variables: 変数
ENV APP_ENV=production  # env: environment（環境）、production: 本番
ENV LOG_LEVEL=info      # log: ログ、level: レベル、info: 情報

# Run the application
# run: 実行する
CMD ["./main"]  # main: メインプログラム 