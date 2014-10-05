#!/usr/bin/env bats

ANDROID_SYMBOLS_DOCKER_REPOSITORY=symbols-help
load $BATS_TEST_DIRNAME/../setup_teardown

@test "symbols help-update" {
	run symbols help-update
	[[ "$status" -eq 0 ]]
	[[ "$output" != "" ]]
}

@test "symbols help-fingerprints" {
	run symbols help-fingerprints
	[[ "$status" -eq 0 ]]
	[[ "$output" != "" ]]
}

@test "symbols help-ls" {
	run symbols help-ls
	[[ "$status" -eq 0 ]]
	[[ "$output" != "" ]]
}

@test "symbols help-fetch" {
	run symbols help-fetch
	[[ "$status" -eq 0 ]]
	[[ "$output" != "" ]]
}

@test "symbols help-pull" {
	run symbols help-pull
	[[ "$status" -eq 0 ]]
	[[ "$output" != "" ]]
}

@test "symbols help-push" {
	run symbols help-push
	[[ "$status" -eq 0 ]]
	[[ "$output" != "" ]]
}

@test "symbols help-hack" {
	run symbols help-hack
	[[ "$status" -eq 0 ]]
	[[ "$output" != "" ]]
}

@test "symbols help-dockerfile" {
	run symbols help-dockerfile
	[[ "$status" -eq 0 ]]
	[[ "$output" != "" ]]
}

@test "symbols help" {
	run symbols help
	[[ "$status" -eq 0 ]]
	[[ "$output" != "" ]]
}

@test "symbols prints help" {
	symbols > >(tee $BATS_TMPDIR/symbols_no_command) || true
	symbols help > >(tee $BATS_TMPDIR/symbols_help)
	diff $BATS_TMPDIR/symbols_no_command $BATS_TMPDIR/symbols_help
}

