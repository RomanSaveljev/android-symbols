#!bash

function teardown {
	docker images -q symbols-bats | xargs docker rmi -f
}

@test "Docker image builds" {
	docker build -t symbols-bats $BATS_TEST_DIRNAME/../..
}
