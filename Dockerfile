FROM docker.io/library/node:22-bookworm-slim AS dashboard_builder

WORKDIR /app/webui
RUN apt-get update \
    && apt-get install -y --no-install-recommends protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*
COPY webui/package.json webui/package-lock.json ./
RUN npm ci
COPY proto /app/proto
COPY webui ./
RUN npm run build

FROM docker.io/library/golang:1.26-alpine AS service_builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o sms-service ./cmd/sms-service

FROM docker.io/library/alpine:3.22

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=service_builder /app/sms-service ./sms-service
COPY --from=dashboard_builder /app/webui/dist /app/dashboard/sms
EXPOSE 50051 8080
CMD ["./sms-service"]
