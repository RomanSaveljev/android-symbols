#!/usr/bin/env bats

function setup {
	. <(docker run --rm $ANDROID_SYMBOLS_DOCKER_REPOSITORY envsetup)
}

function teardown {
	docker ps -q | xargs --no-run-if-empty docker stop
	docker ps -qa | xargs --no-run-if-empty docker rm
	docker images -q symbols-test | sort | uniq | xargs --no-run-if-empty docker rmi -f
}

@test "symbols dockerfile returns template Dockerfile" {
	symbols dockerfile > >(tee $BATS_TMPDIR/symbols_dockerfile)
	diff $BATS_TMPDIR/symbols_dockerfile $BATS_TEST_DIRNAME/../../docker/opt/Dockerfile
}

@test '$STRIP_SYMBOLS_PATH is respected' {
	echo ----- Building with customized Dockerfile
	local tmp=$(mktemp -d -p $BATS_TMPDIR)
	mkdir -p $tmp
	symbols dockerfile >$tmp/Dockerfile
	sed -i 's?^FROM symbols?FROM symbols-bats?' $tmp/Dockerfile
	sed -i 's?^ENV STRIP_SYMBOLS_PATH.*?ENV STRIP_SYMBOLS_PATH 2?' $tmp/Dockerfile
	cat $tmp/Dockerfile
	docker build -t symbols-test $tmp

	echo ----- Sourcing envsetup
	ANDROID_SYMBOLS_DOCKER_REPOSITORY=symbols-test
	. <(docker run --rm symbols-test envsetup)

	echo ----- Updating the layer
	mkdir -p $BATS_TMPDIR/a/b/c/d/e/
	echo 123 >$BATS_TMPDIR/a/b/c/d/e/file
	echo "ro.build.fingerprint=customization" >$BATS_TMPDIR/a/b/build.prop
	symbols update latest < <(tar -czf - -C $BATS_TMPDIR a/b/c/d/e/file a/b/build.prop)

	echo ----- Checking the listing
	run symbols ls latest
	[[ "$status" -eq 0 ]]
	for line in ${lines[@]}
	do
		[[ "$line" == "build.prop" || "$line" == "c/d/e/file" ]]
	done
	[[ "${#lines[@]}" -eq 2 ]]
}
