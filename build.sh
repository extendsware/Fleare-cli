#!/bin/bash

PROJECT_NAME=""
VERSION=""
LAST_COMMIT_ID=$(git rev-parse --short HEAD)
DATE=$(date -u +"%Y-%m-%d")

while IFS=": " read -r key value; do
  key=$(echo "$key" | xargs)
  value=$(echo "$value" | xargs)

  case "$key" in
    project)
      PROJECT_NAME="$value"
      ;;
    version)
      VERSION="$value"
      ;;
  esac
done < <(grep '^[^[:space:]]\+:' build-config.yaml)

# Create version.go file
cat <<EOF > cmd/version.go
package cmd

var (
	ProjectName = "${PROJECT_NAME}"
	Version     = "${VERSION}"
    BuildDate   = "${DATE}"
)
EOF

go-builder build --target all
# go-builder build 