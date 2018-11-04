#!/usr/bin/env bats

load "../support/lib"

type apt >/dev/null 2>&1

exists=$?

@test "availability" {

  run plugin_script "load_input.txt"

  [ "$status" -eq $exists ]
}

@test "unavailability" {
  if [ "$exists" -eq 1 ]; then
    skip
  fi

  alias apt=''

  skip
}

@test "adds package successfully" {
  
  if [ "$exists" -eq 1 ]; then
    skip
  fi

  skip
}

@test "removes package successfully" {
  if [ "$exists" -eq 1 ]; then
    skip
  fi

  skip
}

@test "removes unknown package unsuccessfully" {
  if [ "$exists" -eq 1 ]; then
    skip
  fi

  skip
}

