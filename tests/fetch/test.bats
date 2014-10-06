#!/usr/bin/env bats

ANDROID_SYMBOLS_DOCKER_REPOSITORY=symbols-fetch
load $BATS_TEST_DIRNAME/../setup_teardown

@test "fetch preserves paths" {
	echo "ro.build.fingerprint=preserves-paths" >$BATS_TMPDIR/build.prop
	mkdir -p $BATS_TMPDIR/x/y/z
	echo 123 >$BATS_TMPDIR/x/y/z/file_one
	echo 456 >$BATS_TMPDIR/x/y/file_two
	echo 789 >$BATS_TMPDIR/x/file_three
	echo ABC >$BATS_TMPDIR/file_four
	symbols update latest < <(tar czf - -C $BATS_TMPDIR build.prop x/y/z/file_one x/y/file_two x/file_three file_four)

	local tmp=$(mktemp -d -p $BATS_TMPDIR)
	echo ---- temporary dir $tmp
	symbols fetch latest file_four > >(tar xvf - -C $tmp)
	[[ -e "$tmp/file_four" ]]
	symbols fetch latest x/file_three > >(tar xvf - -C $tmp)
	[[ -e "$tmp/x/file_three" ]]
	symbols fetch latest x/y/file_two > >(tar xvf - -C $tmp)
	[[ -e "$tmp/x/y/file_two" ]]
	symbols fetch latest x/y/z/file_one > >(tar xvf - -C $tmp)
	[[ -e "$tmp/x/y/z/file_one" ]]
}

@test "fetch retrieves all files" {
	echo "ro.build.fingerprint=all-files" >$BATS_TMPDIR/build.prop
	local tmp_one=$(mktemp -p $BATS_TMPDIR)
	local tmp_two=$(mktemp -p $BATS_TMPDIR)
	symbols update latest < <(tar czf - -C $BATS_TMPDIR build.prop $(basename $tmp_one) $(basename $tmp_two))
	local tmp_dest=$(mktemp -p $BATS_TMPDIR -d)
	symbols fetch latest > >(tar xvf - -C $tmp_dest)
	[[ -e $tmp_dest/$(basename $tmp_one) ]]
	[[ -e $tmp_dest/$(basename $tmp_two) ]]
	[[ -e $tmp_dest/build.prop ]]
}

