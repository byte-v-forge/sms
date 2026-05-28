#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SOURCE_ROOT="${SOURCE_ROOT:-$(cd "${ROOT}/../.." && pwd)}"
SMS_PROTO_DIR="${SMS_PROTO_DIR:-${SOURCE_ROOT}/sms/proto}"
COMMON_PROTO_DIR="${COMMON_PROTO_DIR:-${SOURCE_ROOT}/common-lib/proto}"
OUT_DIR="${OUT_DIR:-${ROOT}/src/proto}"
LOCAL_PLUGIN="${ROOT}/node_modules/.bin/protoc-gen-ts_proto"
AGGREGATE_PLUGIN="${SOURCE_ROOT}/webui/node_modules/.bin/protoc-gen-ts_proto"
PLUGIN="${PROTOC_GEN_TS_PROTO:-}"

if [[ -z "${PLUGIN}" ]]; then
  if [[ -x "${LOCAL_PLUGIN}" ]]; then
    PLUGIN="${LOCAL_PLUGIN}"
  elif [[ -x "${AGGREGATE_PLUGIN}" ]]; then
    PLUGIN="${AGGREGATE_PLUGIN}"
  fi
fi

if [[ -z "${PLUGIN}" || ! -x "${PLUGIN}" ]]; then
  printf 'ts-proto plugin not found; run npm install in webui first\n' >&2
  exit 1
fi
SMS_INTERNAL_PROTO="${SMS_PROTO_DIR}/byte/v/forge/sms/internal/v1/sms_internal.proto"
SMS_CONTRACT_PROTO="${COMMON_PROTO_DIR}/byte/v/forge/contracts/sms/v1/sms.proto"
if [[ ! -f "${SMS_INTERNAL_PROTO}" || ! -f "${SMS_CONTRACT_PROTO}" ]]; then
  printf 'sms proto not found under: %s or %s\n' "${SMS_PROTO_DIR}" "${COMMON_PROTO_DIR}" >&2
  exit 1
fi

rm -rf "${OUT_DIR}"
mkdir -p "${OUT_DIR}"

PROTO_INCLUDES=("-I" "${SMS_PROTO_DIR}" "-I" "${COMMON_PROTO_DIR}")
if [[ -d /usr/include/google/protobuf ]]; then
  PROTO_INCLUDES+=("-I" "/usr/include")
fi

protoc "${PROTO_INCLUDES[@]}" \
  --plugin="protoc-gen-ts_proto=${PLUGIN}" \
  --ts_proto_out="${OUT_DIR}" \
  --ts_proto_opt=onlyTypes=true,outputServices=none,esModuleInterop=true,useJsonWireFormat=true,snakeToCamel=false,Mbyte/v/forge/contracts/sms/v1/sms.proto=@byte-v-forge/common-ui/proto/byte/v/forge/contracts/sms/v1/sms \
  "${SMS_INTERNAL_PROTO}"
