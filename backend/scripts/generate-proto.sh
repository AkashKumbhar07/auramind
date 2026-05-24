#!/usr/bin/env bash
set -euo pipefail

BACKEND_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PROTO_DIR="$BACKEND_DIR/framework/grpc/proto"

export PATH="$HOME/go/bin:$PATH"

echo "Generating gRPC code from proto files..."
echo "  Proto dir: $PROTO_DIR"
echo "  Backend dir: $BACKEND_DIR"

protoc \
  --proto_path="$PROTO_DIR" \
  --go_out="$BACKEND_DIR" \
  --go_opt=module=github.com/AkashKumbhar07/auramind/backend \
  --go-grpc_out="$BACKEND_DIR" \
  --go-grpc_opt=module=github.com/AkashKumbhar07/auramind/backend \
  "$PROTO_DIR"/*.proto

echo "Done."
