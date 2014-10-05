#!/usr/bin/env bats

function setup {
	. <(docker run --rm $ANDROID_SYMBOLS_DOCKER_REPOSITORY envsetup)
}

@test "symbols dockerfile returns template Dockerfile" {
	symbols dockerfile > >(tee $BATS_TMPDIR/symbols_dockerfile)
	diff $BATS_TMPDIR/symbols_dockerfile ../../docker/opt/Dockerfile
}
