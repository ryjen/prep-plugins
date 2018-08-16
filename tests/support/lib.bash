#!/usr/bin/env bash

export TEST_ROOT_DIRNAME=$(dirname $BATS_TEST_DIRNAME)

export TEST_PLUGIN=${BATS_TEST_FILENAME/tests/src}
export TEST_PLUGIN=${TEST_PLUGIN/.bats/}

function plugin_script() {

  local input_file="$1"

  shift

  local plugin=${TEST_PLUGIN}

  if [ ! -x "$plugin" ]; then
    echo "Invalid plugin for execute (${plugin})"
    return -1
  fi

  echo "executing ${plugin} on ${TEST_ROOT_DIRNAME}/data/${input_file}"

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

function plugin_build() {

  local plugin=${TEST_PLUGIN}

  if [ ! -x "${plugin}" ]; then
    echo "Invalid plugin for execute (${plugin})"
    return -1
  fi

  ${TEST_ROOT_DIRNAME}/methods/build $@ | ${plugin}
}

function errecho {
  >&2 echo "$@"
}

