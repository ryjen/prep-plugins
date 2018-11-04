#!/usr/bin/env bats

load "../support/lib"

type make > /dev/null 2>&1

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

@test "valid build" {

  local SRCDIR="${TEST_ROOT_DIRNAME}/data/make"

  if [ ! -d $SRCDIR ]; then
    echo "${SRCDIR} not a directory"
    return 1
  fi

  run plugin_build "${SRCDIR}"

  [ "$status" -eq 0 ]
}

@test "invalid build" {
  local SRCDIR="/tmp/invalid"

  run plugin_build "${SRCDIR}"

  [ "$status" -ne 0 ]
}

