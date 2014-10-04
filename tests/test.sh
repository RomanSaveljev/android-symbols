#!/bin/bash

#source $(dirname $0)/setup.source
ALL_TEST_BATS=$(find -type f -name test.bats)
bats "$@" $ALL_TEST_BATS
