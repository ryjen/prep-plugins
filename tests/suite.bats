#!/usr/bin/env bats

@test "git plugin" {
  run bats git
  [ $status -eq 0 ]
}

@test "cmake plugin" {
  run bats cmake
  [ $status -eq 0 ]
}
