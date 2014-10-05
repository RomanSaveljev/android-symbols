#!/usr/bin/env bats

@test "No syntax errors in docker/opt/envsetup.source" {
	bash -n $BATS_TEST_DIRNAME/../../docker/opt/envsetup.source
}
