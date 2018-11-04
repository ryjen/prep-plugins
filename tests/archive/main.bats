#!/usr/bin/env bats

load "../support/lib"

type grater > /dev/null 2>&1

exists=$?

@test "load available" {

  run plugin_script "load_input.txt"

  [ $status -eq $exists ]
}

@test "load unavailable" {  

  if [ "$exists" -ne 0 ]; then
    skip
  fi

  alias grater=''

  run plugin_script "load_input.txt"

  [ "$status" -eq 0 ]
}

@test "resolve successfully" {

  if [ "$exists" -ne 0 ]; then
    skip
  fi

  run plugin_resolve "file://${TEST_ROOT_DIRNAME}/data/archive/test-0.1.0.tar.gz"
  
  [ "$status" -eq 0 ]
}

@test "resolve invalid location" {
  if [ "$exists" -ne 0 ]; then
    skip
  fi

  run plugin_resolve "gobbdlygook"

  [ "$status" -ne 0 ]
}


