#!/bin/bash

export ANDROID_SYMBOLS_DOCKER_REPOSITORY=symbols-bats

function cleanup_docker {
	docker images -q $ANDROID_SYMBOLS_DOCKER_REPOSITORY | xargs docker rmi -f
}
trap cleanup_docker EXIT
docker build -t $ANDROID_SYMBOLS_DOCKER_REPOSITORY $(dirname $0)/..

#source $(dirname $0)/setup.source
ALL_TEST_BATS=$(find -type f -name test.bats)
bats "$@" $ALL_TEST_BATS
