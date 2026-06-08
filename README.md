# sms

独立 SMS provider 集成服务，包含 provider adapter、SMS 契约、取号/验证码生命周期、gRPC API、HTTP dashboard BFF 和自带静态 Web UI。

## 当前实现

- Go module：`github.com/byte-v-forge/sms`
- 运行入口：`cmd/sms-service`
- 公开契约：`proto/byte/v/forge/contracts/sms/v1/sms.proto`
- 内部管理契约：`proto/byte/v/forge/sms/internal/v1/sms_internal.proto`
- gRPC adapter：`internal/adapters/grpc`
- 可选事件 worker：`internal/adapters/eventbus`
- 核心生命周期服务：`internal/app`
- provider 插件注册：`internal/app/provider_plugin.go`
- provider adapter：`internal/providers/fivesim`、`internal/providers/herosms`、`internal/providers/smsbower`
- 自带 UI：`webui/`，由 SMS 服务 HTTP 入口直接托管

## 可选基础设施

- PostgreSQL 可选：配置 `SMS_PG_DSN` 后持久化 provider 配置、订单、验证码引用和事件 outbox；未配置时使用进程内存储。
- Redis 可选：配置 `SMS_REDIS_URL` 后保存激活热状态、验证码 TTL secret 和路由熔断状态；未配置时使用进程内 TTL 存储。
- NATS/JetStream 可选：配置 `SMS_NATS_URL` 且启用 PostgreSQL outbox 后使用异步取号、轮询和取消 worker；未配置时取号/取消走进程内同步流程，hotstream 使用本地 hub。

## 数据流

- `SmsCatalogService.ListSmsPriceOffers` 查询已启用 provider，返回报价、库存、上游 provider 标识和可放入取号请求的 `SmsOfferRef`。
- `SmsCatalogService.RecommendSmsRoutes` 基于最低价优先、库存最低值、价格上下限、provider 过滤和连续失败 TTL 临时禁用名单提供候选排序。
- `AcquireNumber` 只接受精确 provider 参数；启用事件总线和 outbox 时发布 `sms.order.acquire_requested`，否则在请求内直接调用 provider 取号。
- poll worker 或同步 `GetOrder` 查询 provider 状态；收到验证码后先写入 TTL secret store，再把 `SecretRef` 写入订单验证码历史。
- cancel API 对已收到验证码的订单不再向 provider 取消；若取消竞态中已收到验证码，则同步状态后返回当前订单。

## 契约

- 公共 SMS 契约：`proto/byte/v/forge/contracts/sms/v1/sms.proto`
- 公共基础契约：`proto/byte/v/forge/contracts/common/v1/`
- SMS 内部契约：`proto/byte/v/forge/sms/internal/v1/sms_internal.proto`
- provider 私有契约：
  - `proto/byte/v/forge/sms/providers/fivesim/v1/fivesim.proto`
  - `proto/byte/v/forge/sms/providers/herosms/v1/herosms.proto`
  - `proto/byte/v/forge/sms/providers/smsbower/v1/smsbower.proto`

生成代码：

```sh
sh scripts/generate-proto.sh
cd webui && npm install --ignore-scripts && npm run proto
```

## 运行配置

- `SMS_LISTEN_ADDR`：gRPC 监听地址，默认 `:50051`
- `SMS_DASHBOARD_HTTP_ADDR`：dashboard BFF 与静态 UI 监听地址，默认 `:8080`
- `SMS_DASHBOARD_STATIC_DIR`：静态 UI 目录，默认 `/app/dashboard/sms`
- `SMS_PG_DSN`：可选 PostgreSQL 连接
- `SMS_NATS_URL`：可选 NATS/JetStream 事件总线
- `SMS_EVENT_STREAM_NAME`：可选事件 stream 名称
- `SMS_REDIS_URL`：可选 Redis 连接
- `SMS_PROVIDER_HTTP_PROXY`：可选 provider 出口代理，支持 `http`、`https`、`socks5`、`socks5h`

provider API key 通过 SMS 设置页或 `SmsProviderAdminService` 写入。

## Provider API 与 SDK 取舍

- `5sim`：按官方 API 文档的 `guest/countries`、`guest/prices`、`user/buy/activation`、`user/check`、`user/finish`、`user/cancel` 和 `user/profile` 端点实现；官方未提供稳定 Go SDK，本仓保留轻量 HTTP adapter。
- `smsbower`：按官方 `handler_api.php` 协议的 `getNumberV2`、`getStatus`、`setStatus`、`getBalance`、`getServicesList`、`getCountries` 和 `getPricesV3` action 实现；官方未提供稳定 Go SDK，本仓复用内部 `handlerapi` adapter。
- `herosms`：按官方 SMS-Activate 兼容协议实现取号、查码、改状态和余额查询；报价目录使用 HeroSMS 自有 API 端点；官方未提供稳定 Go SDK，本仓复用内部 `handlerapi` adapter 并只保留必要 OpenAPI 查询客户端。

当前未引入第三方非官方 SDK，避免把 provider 状态码、错误码和凭据处理交给维护状态不明确的依赖。

## 镜像

在目标构建环境中以本仓库为 build context 构建：

```sh
docker build -t sms-service .
```
