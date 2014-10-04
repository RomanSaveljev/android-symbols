#!bash

function setup {
	docker build -t symbols-bats $BATS_TEST_DIRNAME/../..
}

function teardown {
	docker images -q symbols-bats | xargs docker rmi -f
}

@test "envsetup prints script" {
	skip
	run docker run --rm symbols-bats envsetup
	[ "$output" != "" ]
}

@test "symbols function prints help" {
	source <(docker run --rm symbols-bats envsetup)
	run symbols
	echo $output
	[ "$status" -eq 0 ]
	[ "$output" != "" ]
}

