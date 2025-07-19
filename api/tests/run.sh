#!/usr/bin/bash
cd "$(dirname "$0")"
docker compose --env-file ../.env up -d --build
go clean -testcache
go test ./... -v -p 4
# docker compose down