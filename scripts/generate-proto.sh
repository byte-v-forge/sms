#!/usr/bin/env sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
SOURCE_ROOT="${SOURCE_ROOT:-$(CDPATH= cd -- "$ROOT/.." && pwd)}"
COMMON_PROTO_DIR="${COMMON_PROTO_DIR:-$SOURCE_ROOT/common-lib/proto}"
PATH="$(go env GOPATH)/bin:$PATH"

if [ ! -f "$COMMON_PROTO_DIR/byte/v/forge/contracts/sms/v1/sms.proto" ]; then
  printf 'common sms contract not found under: %s\n' "$COMMON_PROTO_DIR" >&2
  exit 1
fi

rm -rf "$ROOT/gen"
mkdir -p "$ROOT/gen/go"

protoc -I "$ROOT/proto" -I "$COMMON_PROTO_DIR" \
  --go_out="$ROOT" \
  --go_opt=module=github.com/byte-v-forge/sms \
  --go-grpc_out="$ROOT" \
  --go-grpc_opt=module=github.com/byte-v-forge/sms \
  "$ROOT/proto/byte/v/forge/sms/internal/v1/sms_internal.proto" \
  "$ROOT/proto/byte/v/forge/sms/providers/fivesim/v1/fivesim.proto" \
  "$ROOT/proto/byte/v/forge/sms/providers/herosms/v1/herosms.proto" \
  "$ROOT/proto/byte/v/forge/sms/providers/smsbower/v1/smsbower.proto"

gofmt -w "$ROOT/gen/go"
