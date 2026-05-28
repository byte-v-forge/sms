FROM docker.m.daocloud.io/library/node:22-bookworm-slim AS dashboard_remote_builder

WORKDIR /sms/webui
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources     && apt-get update     && apt-get install -y --no-install-recommends libprotobuf-dev protobuf-compiler     && rm -rf /var/lib/apt/lists/*
COPY common-lib/ui /common-lib/ui
COPY common-lib/proto /common-lib/proto
COPY common-lib/scripts /common-lib/scripts
COPY sms/proto /sms/proto
COPY sms/webui ./
RUN npm ci && SOURCE_ROOT=/ npm run build

FROM docker.m.daocloud.io/library/golang:1.26-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache git ca-certificates

COPY common-lib /common-lib
COPY sms/go.mod sms/go.sum ./
RUN go mod edit -replace github.com/byte-v-forge/common-lib=/common-lib \
    && go mod download

COPY sms .
RUN CGO_ENABLED=0 GOOS=linux go build -o sms-service ./cmd/sms-service

FROM docker.m.daocloud.io/library/alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/sms-service .
COPY --from=dashboard_remote_builder /sms/webui/dist /app/dashboard/sms
EXPOSE 50051 8080
CMD ["./sms-service"]
