# 贡献指南

## 边界

本仓只承载 SMS 业务能力、provider adapter、本仓 proto、SMS 服务实现和自带 Web UI。

以下内容不进入本仓：

- 其他业务域代码；
- 真实 provider 凭据、真实手机号或真实验证码；
- 检查报告、临时二进制和其他构建产物。

## 开发流程

1. 公开 SMS 契约、公共基础契约、业务内部模型和 provider 私有 shape 均在本仓 `proto/` 修改。
2. provider 私有 shape 放在 `proto/byte/v/forge/sms/providers/<provider>/v1/`。
3. 业务内部模型放在 `proto/byte/v/forge/sms/internal/v1/`。
4. 外部 provider 调用必须设置超时，并按 provider 官方文档实现状态和错误映射。
5. Web UI 必须作为本服务自带静态应用，不使用 Module Federation 或外部 common-ui 包。

## 验证

```sh
sh scripts/generate-proto.sh
go list ./...
cd webui && npm run lint
```

`gen/` 只承载本仓 proto 生成物。
