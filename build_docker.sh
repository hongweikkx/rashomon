#!/usr/bin/env bash

set -exo pipefail

while getopts ":v:p:" opt; do
    case "$opt" in
    v) VERN="${OPTARG}" ;;
    p) PUSH="${OPTARG}" ;;
    *)
    esac
done

if [[ -z "$VERN" ]];then
  echo ERROR: VERSION is nil
  exit 1
fi

make vendor

IMAGE_NAME=registry.cn-hangzhou.aliyuncs.com/hongweigao/rashomon:${VERN}

docker build -t ${IMAGE_NAME} -f docker/Dockerfile .


if [[ "$PUSH" = "true" ]]; then
  docker push ${IMAGE_NAME}
fi

