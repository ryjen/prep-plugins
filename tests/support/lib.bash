#!/usr/bin/env bash

TEST_ROOT_DIRNAME=$(dirname $BATS_TEST_DIRNAME)

TEST_PLUGIN=${BATS_TEST_FILENAME/tests/src}

function plugin_script() {

  local input_file="$1"

  shift

  local plugin=${TEST_PLUGIN}

  if [ ! -x "$plugin" ]; then
    echo "Invalid plugin for execute (${plugin})"
    return -1
  fi

  ${plugin} < ${TEST_ROOT_DIRNAME}/data/${input_file}
}

function plugin_resolve() {

  local plugin=${TEST_PLUGIN}

  if [ ! -x "$plugin" ]; then
    echo "Invalid plugin for execute (${plugin})"
    return -1
  fi

  ${TEST_ROOT_DIRNAME}/methods/resolver $@ | ${plugin}
}

function errecho {
    >&2 echo "$@"
}

