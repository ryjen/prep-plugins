#!/usr/bin/env bats

load "../support/lib"

type cmake > /dev/null 2>&1

exists=$?

@test "available" {
  
  run plugin_script "load_input.txt"

  [ "$status" -eq $exists ]
}

@test "unavailable" {
  
  if [ "$exists" -ne 0 ]; then
    skip
  fi

  run plugin_script "load_input.txt"

  [ "$status" -eq 0 ]
}


