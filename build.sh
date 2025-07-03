#!/usr/bin/env sh
set -e
set -x

export PACKAGE_NAME="filament-sync-tool"
export MODULE_PATH="cli"

# disable CGO for cross-compiling
export CGO_ENABLED=0

# compile for all architectures
GOOS=linux   GOARCH=amd64   go build  -o release/linux/amd64/${PACKAGE_NAME}       ./${MODULE_PATH}
GOOS=linux   GOARCH=arm64   go build  -o release/linux/arm64/${PACKAGE_NAME}       ./${MODULE_PATH}
GOOS=windows GOARCH=amd64   go build  -o release/windows/amd64/${PACKAGE_NAME}.exe ./${MODULE_PATH}
GOOS=darwin  GOARCH=amd64   go build  -o release/darwin/amd64/${PACKAGE_NAME}      ./${MODULE_PATH}
GOOS=darwin  GOARCH=arm64   go build  -o release/darwin/arm64/${PACKAGE_NAME}      ./${MODULE_PATH}

# tar binary files prior to upload
tar -cvzf release/${PACKAGE_NAME}_linux_amd64.tar.gz   -C release/linux/amd64   ./${PACKAGE_NAME}
tar -cvzf release/${PACKAGE_NAME}_linux_arm64.tar.gz   -C release/linux/arm64   ./${PACKAGE_NAME}
tar -cvzf release/${PACKAGE_NAME}_windows_amd64.tar.gz -C release/windows/amd64 ./${PACKAGE_NAME}.exe
tar -cvzf release/${PACKAGE_NAME}_darwin_amd64.tar.gz  -C release/darwin/amd64  ./${PACKAGE_NAME}
tar -cvzf release/${PACKAGE_NAME}_darwin_arm64.tar.gz  -C release/darwin/arm64  ./${PACKAGE_NAME}

# generate shas for tar files
sha256sum release/*.tar.gz > release/${PACKAGE_NAME}_checksums.txt


## Publish local build
# git tag X.y.Z
# git push --tags
# (gh release delete  <tag> --cleanup-tag)
# gh release create <tag> ./release/*.gz
