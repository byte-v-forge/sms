# AGENTS.md

- 本仓库只承载 SMS 业务能力、provider adapter、业务内部契约和 SMS 服务实现。
- 公共契约来自 `common-lib/proto/byte/v/forge/contracts/sms/v1/`；provider 私有 shape 留在本仓库内部 proto。
- 后端优先使用 Go，按 Clean Code、DI 和面向抽象设计组织代码。
- `gen/` 只承载本仓内部/provider proto 生成物；公开 Go 契约生成物来自 `common-lib/gen/go/byte/v/forge/contracts/sms/v1/`；检查报告、临时二进制和其他构建产物不提交。
- proto 变更后必须运行生成命令、格式化和 Go 检查。
