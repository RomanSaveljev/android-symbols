#!/usr/bin/env bats

function teardown {
	docker images -q symbols-build | xargs --no-run-if-empty docker rmi --no-prune
	docker images -q symbols-customized | xargs --no-run-if-empty docker rmi -f
}

@test "Docker image builds" {
	docker build -t symbols-build $BATS_TEST_DIRNAME/../..
}

@test "Customized image builds" {
	local tmp=$(mktemp -p $BATS_TMPDIR)
	docker build -t symbols-customized .
}
