#!/bin/bash

# setup
folder="generated"
package="generated"
oapiFile="openapi.yaml"

# cleanup
rm -rf ./handler/$folder
mkdir -p ./handler/$folder

# codegen
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
oapi-codegen -generate types -o handler/$folder/$package-types.gen.go -package $package spec/$oapiFile
oapi-codegen -generate "chi-server" -o handler/$folder/$package.gen.go -package $package spec/$oapiFile

# go
go mod tidy
