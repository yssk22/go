#!/bin/sh

VERSION=$(git describe --tags 2>/dev/null || git describe --contains --all HEAD)
BUILD_REV=$(git describe --always)

echo "VERSION=${VERSION}"
echo "BUILD REV=${BUILD_REV}"

CURRENT_DIR=$(dirname $0)

TOOL_DIR="${CURRENT_DIR}/../tools/cmd"

DIST_DIR=${CURRENT_DIR}/../dist
mkdir -p $DIST_DIR

MAIN_FILES=$(grep -R "func main()" ${TOOL_DIR} | awk -F: '{print $1}')

for MAIN in ${MAIN_FILES}
do
    MAIN_DIR=$(dirname $MAIN)
    PKG_NAME=$(basename $MAIN_DIR)
    gox -output "dist/${PKG_NAME}.{{.OS}}.{{.Arch}}" \
        -ldflags "-X main.BuildRev $BUILD_REV -X main.Version $VERSION" \
        -os="darwin" -arch="amd64" \
        $MAIN_DIR
done
