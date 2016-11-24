#!/bin/bash

CURRENT_DIR=$(dirname $0)
DIST_DIR="${CURRENT_DIR}/../dist"

GHR=$(which ghr)

if [[ "${GHR}" = "" ]]; then
    if [[ ! -f "./ghr" ]]; then
        wget -O ghr.zip https://github.com/tcnksm/ghr/releases/download/v0.5.3/ghr_v0.5.3_linux_amd64.zip 2>/dev/null
        unzip ghr.zip
    fi
    GHR="./ghr"
fi

VERSION=$(git describe --tags 2>/dev/null || git describe --contains --all HEAD)

${GHR} -u $CIRCLE_PROJECT_USERNAME -r $CIRCLE_PROJECT_REPONAME --replace ${VERSION} ${DIST_DIR}