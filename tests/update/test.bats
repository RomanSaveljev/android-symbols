#!/usr/bin/env bats

ANDROID_SYMBOLS_DOCKER_REPOSITORY=symbols-update

function setup {
	docker build -t symbols-bats $BATS_TEST_DIRNAME/../..
	docker build -t $ANDROID_SYMBOLS_DOCKER_REPOSITORY .
	. <(docker run --rm symbols-bats envsetup)
}

function teardown {
	docker ps -q | xargs --no-run-if-empty docker stop
	docker ps -qa | xargs --no-run-if-empty docker rm
	docker images -q $ANDROID_SYMBOLS_DOCKER_REPOSITORY | xargs --no-run-if-empty docker rmi -f
	docker images -q symbols-bats | xargs --no-run-if-empty docker rmi -f
}

@test "First batch of symbols is stored" {
	skip
	local files_path=$(mktemp -d -p $BATS_TMPDIR)
	mkdir -p $files_path
	echo "ro.build.fingerprint=my-great-fingerprint" >$files_path/build.prop
	dd if=/dev/zero of=$files_path/first bs=1M count=12
	dd if=/dev/zero of=$files_path/second bs=1M count=1
	dd if=/dev/zero of=$files_path/third bs=1M count=3
	symbols update latest < <(tar czvf - -C $files_path .)
}

@test "Update fails, if called without base layer id" {
	local tmp=$(mktemp -p $BATS_TMPDIR)
	local build_prop=$BATS_TMPDIR/build.prop
	echo "ro.build.fingerprint=my-great-fingerprint" >$build_prop
	dd if=/dev/zero of=$tmp bs=1M count=1
	run symbols update < <(tar czvf - -C $BATS_TMPDIR $(basename $tmp) $(basename $build_prop))
	[[ "$status" -ne 0 ]]
}

@test "Update fails, if build.prop is not included" {
	local tmp=$(mktemp -p $BATS_TMPDIR)
	dd if=/dev/zero of=$tmp bs=1M count=1
	run symbols update latest < <(tar czvf - -C $BATS_TMPDIR $(basename $build_prop))
	[[ "$status" -ne 0 ]]
}

@test "Update fails, if STDIN is closed" {
	run symbols update latest </dev/null
	[[ "$status" -ne 0 ]]
}

@test "Update fails, if wrong base layer is selected" {
}
