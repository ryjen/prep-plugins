#!/usr/bin/env bats

@test "git plugin" {
  run bats ${BATS_TEST_DIRNAME}/git
  [ $status -eq 0 ]
}

@test "cmake plugin" {
  run bats ${BATS_TEST_DIRNAME}/cmake
  [ $status -eq 0 ]
}

@test "autotools plugin" {
  run bats ${BATS_TEST_DIRNAME}/autotools
  [ $status -eq 0 ]
}

