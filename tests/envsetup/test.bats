#!/usr/bin/env bats

@test "envsetup prints script" {
	run docker run --rm symbols-bats envsetup
	[[ "$status" -eq 0 ]]
	[[ "$output" != "" ]]
	local tmp=$(mktemp -p $BATS_TMPDIR)
	docker run --rm symbols-bats envsetup >$tmp
	diff $tmp $BATS_TEST_DIRNAME/../../docker/opt/envsetup.source
}

@test "symbols without command returns error" {
	source <(docker run --rm symbols-bats envsetup)
	run symbols
	[[ "$status" -ne 0 ]]
}

