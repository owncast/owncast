#!/bin/bash

# go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

# setup
package="generated"
folderPath="webserver/handlers/generated"
specPath="openapi.yaml"

# validate scripts are installed
if ! command -v redocly &>/dev/null; then
	echo "Please install \`redocly cli\` before running this script: npm install -g @redocly/cli"
	exit 1
fi

# validate schema
npx redocly lint $specPath
if [ $? -ne 0 ]; then
	echo "Open API specification is not valid"
	exit 1
fi

# cleanup
rm -r $folderPath
mkdir -p $folderPath

# codegen
go tool oapi-codegen -generate types -o $folderPath/$package-types.gen.go -package $package $specPath
go tool oapi-codegen -generate "chi-server" -o $folderPath/$package.gen.go -package $package $specPath

# go
go mod tidy
