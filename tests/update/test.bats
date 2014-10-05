#!/usr/bin/env bats

ANDROID_SYMBOLS_DOCKER_REPOSITORY=symbols-update
load $BATS_TEST_DIRNAME/../setup_teardown

@test "First batch of symbols is stored" {
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
	echo "ro.build.fingerprint=my-great-fingerprint" >$BATS_TMPDIR/build.prop
	run symbols update asdasdasd < <(tar czvf - -C $files_path $BATS_TMPDIR/build.prop)
	[[ "$status" -ne 0 ]]
}

@test "Update another layer" {
	echo "ro.build.fingerprint=first" >$BATS_TMPDIR/build.prop
	symbols update latest < <(tar czvf - -C $BATS_TMPDIR build.prop)
	echo "ro.build.fingerprint=second" >$BATS_TMPDIR/build.prop
	symbols update first < <(tar czvf - -C $BATS_TMPDIR build.prop)
	echo "ro.build.fingerprint=third" >$BATS_TMPDIR/build.prop
	symbols update first < <(tar czvf - -C $BATS_TMPDIR build.prop)
}

@test "Update same layer" {
	echo "ro.build.fingerprint=same" >$BATS_TMPDIR/build.prop
	dd if=/dev/zero of=$BATS_TMPDIR/first bs=100 count=1
	symbols update latest < <(tar czvf - -C $BATS_TMPDIR build.prop first)
	symbols update same < <(tar czvf - -C $BATS_TMPDIR build.prop first)
	symbols update same < <(tar czvf - -C $BATS_TMPDIR build.prop first)
}

@test "Update deletes missing files" {
	echo "ro.build.fingerprint=delete-missing" >$BATS_TMPDIR/build.prop
	dd if=/dev/zero of=$BATS_TMPDIR/first bs=100 count=1
	symbols update latest < <(tar czf - -C $BATS_TMPDIR build.prop first)
	dd if=/dev/zero of=$BATS_TMPDIR/second bs=100 count=1
	symbols update delete-missing < <(tar czf - -C $BATS_TMPDIR build.prop second)
	run symbols ls latest
	[[ "$status" -eq 0 ]]
	for line in ${lines[@]}
	do
		echo filename $line
		[[ "$line" != "first" ]]
		[[ "$line" == "build.prop" || "$line" == "second" ]]
	done
	[[ "${#lines[@]}" -eq 2 ]]
}

@test "Update stores empty files" {
	echo "ro.build.fingerprint=empty" >$BATS_TMPDIR/build.prop
	echo >$BATS_TMPDIR/empty
	symbols update latest < <(tar czvf - -C $BATS_TMPDIR build.prop empty)
	run symbols ls latest
	[[ "$status" -eq 0 ]]
	for line in ${lines[@]}
	do
		echo ----- $line
		[[ "$line" == "build.prop" || "$line" == "empty" ]]
	done
	[[ "${#lines[@]}" -eq 2 ]]
}
