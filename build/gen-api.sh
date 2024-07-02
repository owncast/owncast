#!/bin/bash

# go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

# setup
package="generated"
folderPath="webserver/handlers/generated"
specPath="openapi.yaml"

# validate scripts are installed
if ! command -v swagger-cli &>/dev/null; then
	echo "Please install \`swagger-cli\` before running this script"
	exit 1
fi

if ! command -v oapi-codegen &>/dev/null; then
	echo "Please install \`oapi-codegen\` before running this script"
	echo "Hint: run \`go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest\` to install"
	exit 1
fi

# validate schema
swagger-cli validate $specPath
if [ $? -ne 0 ]; then
	echo "Open API specification is not valid"
	exit 1
fi

# cleanup
rm -r $folderPath
mkdir -p $folderPath

# codegen
oapi-codegen -generate types -o $folderPath/$package-types.gen.go -package $package $specPath
oapi-codegen -generate "chi-server" -o $folderPath/$package.gen.go -package $package $specPath

# go
go mod tidy
