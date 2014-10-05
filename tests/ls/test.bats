#!/usr/bin/env bats

ANDROID_SYMBOLS_DOCKER_REPOSITORY=symbols-ls
load $BATS_TEST_DIRNAME/../setup_teardown

@test "ls works" {
	local files_path=$(mktemp -d -p $BATS_TMPDIR)
	mkdir -p $files_path
	echo "ro.build.fingerprint=my-great-fingerprint" >$files_path/build.prop
	dd if=/dev/zero of=$files_path/first bs=1M count=12
	dd if=/dev/zero of=$files_path/second bs=1M count=1
	dd if=/dev/zero of=$files_path/third bs=1M count=3
	symbols update latest < <(tar czvf - -C $files_path .)
	run symbols ls latest
	[[ "$status" -eq 0 ]]
	[[ "${#lines[@]}" -eq 4 ]]
	for name in $output
	do
		echo $name
		[[ "$name" == "first" || \
		"$name" == "second" || \
		"$name" == "third" || \
		"$name" == "build.prop" ]]
	done
}

@test "ls without layer id returns error" {
	run symbols ls
	[[ "$status" -ne 0 ]]
}

@test "ls for wrong layer id returns error" {
	run symbols ls dsfsdfsdf
	[[ "$status" -ne 0 ]]
}

@test "ls prints everything recursively" {
	echo "ro.build.fingerprint=recursive" >$BATS_TMPDIR/build.prop
	mkdir -p $BATS_TMPDIR/a/b/c/d/e/f/g
	touch $BATS_TMPDIR/a/b/c/d/e/f/g/file
	symbols update latest < <(tar czvf - -C $BATS_TMPDIR build.prop a/b/c/d/e/f/g/file)
	run symbols ls recursive
	[[ "$status" -eq 0 ]]
	[[ "${#lines[@]}" -eq 2 ]]
	for line in ${lines[@]}
	do
		[[ "$line" == "build.prop" || $line == "a/b/c/d/e/f/g/file" ]]
	done
}
