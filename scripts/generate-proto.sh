#!/usr/bin/env sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
PATH="$(go env GOPATH)/bin:$PATH"

rm -rf "$ROOT/gen"
mkdir -p "$ROOT/gen/go"

protoc -I "$ROOT/proto" \
  --go_out="$ROOT" \
  --go_opt=module=github.com/byte-v-forge/sms \
  --go-grpc_out="$ROOT" \
  --go-grpc_opt=module=github.com/byte-v-forge/sms \
  "$ROOT/proto/byte/v/forge/contracts/common/v1/common.proto" \
  "$ROOT/proto/byte/v/forge/contracts/common/v1/eventbus.proto" \
  "$ROOT/proto/byte/v/forge/contracts/observability/v1/hotstream.proto" \
  "$ROOT/proto/byte/v/forge/contracts/sms/v1/sms.proto" \
  "$ROOT/proto/byte/v/forge/sms/internal/v1/sms_internal.proto" \
  "$ROOT/proto/byte/v/forge/sms/providers/fivesim/v1/fivesim.proto" \
  "$ROOT/proto/byte/v/forge/sms/providers/herosms/v1/herosms.proto" \
  "$ROOT/proto/byte/v/forge/sms/providers/smsbower/v1/smsbower.proto"

gofmt -w "$ROOT/gen/go"
