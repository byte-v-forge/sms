# sms

`sms` 是独立 SMS provider 集成服务，负责取号、验证码生命周期、provider 管理、报价推荐和自带 dashboard。

## 核心能力

- 聚合多个 SMS provider，提供报价、库存、余额、取号、查码、取消和完成等能力。
- 基于 provider capability 与 route recommendation 选择可用号码渠道。
- 管理订单状态、验证码 TTL secret、provider 配置和失败熔断状态。
- 提供 gRPC API、HTTP dashboard BFF 和独立静态 Web UI。
- PostgreSQL、Redis、NATS/JetStream 均为可选增强；未配置时保持 standalone 可启动。

## 使用方式

业务服务通过 SMS 契约申请号码、查询订单或消费验证码事件；provider API key 通过设置页或管理 API 写入，不写入业务仓。验证码原文只进入服务自有 TTL secret store，对外返回引用。

## 入口

- 服务入口：`cmd/sms-service`
- 公开契约：`proto/byte/v/forge/contracts/sms/v1/`
- 内部/provider 契约：`proto/byte/v/forge/sms/`
- Provider adapter：`internal/providers/`
- Dashboard：`webui/`

## 常用检查

```sh
sh scripts/generate-proto.sh
(cd webui && npm run proto)
git diff --check
```
