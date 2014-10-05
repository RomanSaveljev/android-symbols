#!/bin/bash

function cleanup_docker {
	docker images -q symbols-bats | xargs docker rmi -f
}
trap cleanup_docker EXIT
docker build -t symbols-bats $(dirname $0)/..

#source $(dirname $0)/setup.source
ALL_TEST_BATS=$(find -type f -name test.bats)
bats "$@" $ALL_TEST_BATS
