#!/usr/bin/env bash
set -eo pipefail

OPTS="--proto_path=flux/fluxpb --go_out=paths=source_relative:flux/fluxpb --go-grpc_out=paths=source_relative:flux/fluxpb"
protoc ${OPTS} flux/fluxpb/*.proto