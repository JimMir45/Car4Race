# 多阶段构建

# 阶段1: 构建前端
FROM node:22-alpine AS web-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# 阶段2: 构建后端
FROM golang:1.25-alpine AS server-builder
WORKDIR /app/server
RUN apk add --no-cache gcc musl-dev
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=1 go build -o /app/api cmd/api/main.go

# 阶段3: 最终镜像
FROM alpine:latest
WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata sqlite

# 复制构建产物
COPY --from=server-builder /app/api /app/api
COPY --from=web-builder /app/web/dist /app/web/dist

# 创建数据目录
RUN mkdir -p /app/data

# 环境变量
ENV ENV=production
ENV PORT=8080
ENV DB_PATH=/app/data/car4race.db

EXPOSE 8080

CMD ["/app/api"]
