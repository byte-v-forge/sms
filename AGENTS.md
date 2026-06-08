# AGENTS.md

- 本仓库只承载 SMS 业务能力、provider adapter、业务内部契约和 SMS 服务实现。
- SMS 服务可独立于 byte-v-forge 聚合目录构建、运行和部署；不得依赖 `common-lib` 源码、生成物或 UI 包。
- 服务公开/内部契约真源均在本仓库 `proto/`；provider 私有 shape 留在本仓库内部 proto。
- 后端优先使用 Go，按 Clean Code、DI 和面向抽象设计组织代码。
- `gen/` 只承载本仓 proto 生成物；检查报告、临时二进制和其他构建产物不提交。
- PG、Redis、NATS/MQ 均为可选运行能力；缺省运行路径必须使用服务内抽象或内存实现保持可启动。
- Web UI 作为本服务自带静态应用发布；不得使用 Module Federation 或 `@byte-v-forge/common-ui`。
- Linter 检查必须达到 0 error / 0 warning；禁止通过修改或放宽 linter 配置、降低规则级别、删除规则、添加 ignore/disable/nolint/ts-ignore/eslint-disable/biome-ignore/prettier-ignore 等方式绕过问题，只能按 linter 规则修复源码、类型、格式或依赖边界。
- proto 变更后必须运行生成命令、格式化和 Go 检查。
