#! /bin/bash

set -eou pipefail

default_tag="alpine"

DOCKER_REDIS_PORT="${DOCKER_REDIS_PORT:-16379}"

if test $# -eq 0; then
  echo "please specify either 'start' or 'test'"
  exit 1
fi

case "$1" in
  "start")
    docker run --rm -d \
      -p ${DOCKER_REDIS_PORT}:6379 \
      --name redis \
      redis:${2:-$default_tag}
    ;;
  "test")
    TF_ACC=1 \
    REDISDB_HOSTNAME=127.0.0.1 \
    REDISDB_PORT=${DOCKER_REDIS_PORT} \
    REDISDB_DATABASE=1 \
    go test -v -cover -count 1 ./internal/provider
    ;;
  "stop")
    docker stop redis
    ;;
  "update")
    docker pull redis:${2:-$default_tag}
    ;;
  *)
    echo "unrecognized command"
    exit 1
    ;;
esac