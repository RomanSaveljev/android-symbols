#!/usr/bin/env bats

function setup {
	docker run --rm $ANDROID_SYMBOLS_DOCKER_REPOSITORY envsetup
	false
	#source <(docker run --rm $ANDROID_SYMBOLS_DOCKER_REPOSITORY envsetup)
}

@test "symbols dockerfile returns template Dockerfile" {
	[[ "$ANDROID_SYMBOLS_DOCKER_REPOSITORY" == "symbols-bats" ]]
	#run symbols dockerfile
	#diff $BATS_TMPDIR/symbols_dockerfile ../../docker/opt/Dockerfile
}
