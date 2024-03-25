#!/bin/sh

mkdir -p webserver/handlers
oapi-codegen -generate "types" -o webserver/handlers/openapi_server-types.gen.go -package handlers -response-type-suffix responses openapi.yaml
oapi-codegen -generate "chi-server" -o webserver/handlers/openapi_server.gen.go -package handlers -response-type-suffix responses openapi.yaml
