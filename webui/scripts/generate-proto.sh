#!/usr/bin/env sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
SMS_PROTO_DIR="${SMS_PROTO_DIR:-${ROOT}/../proto}"
OUT_DIR="${OUT_DIR:-${ROOT}/src/proto}"
LOCAL_PLUGIN="${ROOT}/node_modules/.bin/protoc-gen-ts_proto"
PLUGIN="${PROTOC_GEN_TS_PROTO:-}"

if [ -z "${PLUGIN}" ] && [ -x "${LOCAL_PLUGIN}" ]; then
  PLUGIN="${LOCAL_PLUGIN}"
fi

if [ -z "${PLUGIN}" ] || [ ! -x "${PLUGIN}" ]; then
  printf 'ts-proto plugin not found; run npm install in webui first\n' >&2
  exit 1
fi

COMMON_PROTO="${SMS_PROTO_DIR}/byte/v/forge/contracts/common/v1/common.proto"
SMS_CONTRACT_PROTO="${SMS_PROTO_DIR}/byte/v/forge/contracts/sms/v1/sms.proto"
SMS_INTERNAL_PROTO="${SMS_PROTO_DIR}/byte/v/forge/sms/internal/v1/sms_internal.proto"
if [ ! -f "${COMMON_PROTO}" ] || [ ! -f "${SMS_CONTRACT_PROTO}" ] || [ ! -f "${SMS_INTERNAL_PROTO}" ]; then
  printf 'sms proto not found under: %s\n' "${SMS_PROTO_DIR}" >&2
  exit 1
fi

rm -rf "${OUT_DIR}"
mkdir -p "${OUT_DIR}"

if [ -d /usr/include/google/protobuf ]; then
  set -- -I "${SMS_PROTO_DIR}" -I /usr/include
else
  set -- -I "${SMS_PROTO_DIR}"
fi

protoc "$@" \
  --plugin="protoc-gen-ts_proto=${PLUGIN}" \
  --ts_proto_out="${OUT_DIR}" \
  --ts_proto_opt=onlyTypes=true,outputServices=none,esModuleInterop=true,useJsonWireFormat=true,snakeToCamel=false \
  "${COMMON_PROTO}" \
  "${SMS_CONTRACT_PROTO}" \
  "${SMS_INTERNAL_PROTO}"
