# sms

SMS provider 集成服务。

本仓库负责 SMS provider adapter、provider 私有契约、激活生命周期、公共 SMS gRPC 实现和 SMS dashboard 模块。

## 当前实现

- Go module：`github.com/byte-v-forge/sms`
- 运行入口：`cmd/sms-service`
- 公共契约 adapter：`internal/adapters/grpc`
- 平台事件 worker：`internal/adapters/eventbus`
- 核心生命周期服务：`internal/app`
- provider 配置表：`sms_provider_configs`
- 轻量订单表：`sms_orders`
- 验证码历史表：`sms_order_codes`
- 平台事件 outbox：`sms_platform_event_outbox`
- 激活热状态：Redis，key 前缀 `sms:order`
- 路由临时熔断状态：Redis，key 前缀 `sms:route_health`
- provider 插件注册：`internal/app/provider_plugin.go`
- provider adapter：`internal/providers/fivesim`、`internal/providers/herosms`、`internal/providers/smsbower`

## 数据流

- `SmsCatalogService.ListSmsPriceOffers` 直接查询已启用 provider，返回报价、库存、上游 provider 标识和可直接取号的 `SmsNumberAcquireParams`。
- `SmsCatalogService.RecommendSmsRoutes` 基于最低价优先、库存最低值、价格上限、provider 过滤和业务方配置的连续失败 TTL 临时禁用名单提供通用候选排序；使用方仍负责最终选择具体报价和上游 provider。
- SMS 服务不在 `AcquireNumber` 中隐式承载自动 fallback 或 route profile。
- `AcquireNumber` 只接受精确 provider 参数，创建本地激活并通过 outbox 发布 `sms.activation.acquire_requested`。
- acquire worker 消费事件后调用对应 provider 精确取号。
- poll worker 查询 provider 状态；收到验证码后先写入 `sms_order_codes`，再通过 outbox 发布公共 SMS 事件。
- cancel API 写入取消意图并发布 `sms.activation.cancel_requested`；cancel worker 执行 provider cancel。
- 已收到验证码的激活不会再向 provider 取消；若取消竞态中已收到验证码，则同步状态后返回当前激活。

## 契约

- 公共 SMS 契约：`common-lib/proto/byte/v/forge/contracts/sms/v1/`
- SMS 内部契约：`proto/byte/v/forge/sms/internal/v1/sms_internal.proto`
- provider 私有契约：
  - `proto/byte/v/forge/sms/providers/fivesim/v1/fivesim.proto`
  - `proto/byte/v/forge/sms/providers/herosms/v1/herosms.proto`
  - `proto/byte/v/forge/sms/providers/smsbower/v1/smsbower.proto`

生成代码：

```sh
sh scripts/generate-proto.sh
sh webui/scripts/generate-proto.sh
```

## 运行配置

- `SMS_LISTEN_ADDR`：gRPC 监听地址，默认 `:50051`
- `SMS_DASHBOARD_HTTP_ADDR`：dashboard BFF 监听地址，默认 `:8080`
- `SMS_PG_DSN` 或 `PG_DSN`：PostgreSQL 连接
- `PLATFORM_NATS_URL`：平台事件总线
- `PLATFORM_REDIS_URL`：激活热状态和路由临时熔断状态 Redis
- `SMS_PROVIDER_HTTP_PROXY`：可选 provider 出口代理，支持 `http`、`https`、`socks5`、`socks5h`

provider API key 通过 SMS 设置页或 `SmsProviderAdminService` 写入数据库。
