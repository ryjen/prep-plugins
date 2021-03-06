#!/usr/bin/env bats

load "../support/lib"

run type git

exists=$status

function resolve_success() {

  readarray -t lines <<<$output

  local values=()
  local valid=0
  
  for value in ${lines[@]}
  do
    command=($value)
    case ${command[0]} in
      ECHO*)
        values+=($value)
        ;;
      RETURN*)
        valid=1
        values+=($value)
        ;;
      *)
        ;;
    esac
  done

  [ "$valid" -eq 1 ]

  [[ "${values[0]}" =~ ^ECHO\ Cloning\ .+$  ]]
  [[ "${values[1]}" =~ ^ECHO\ Updating\ submodules$ ]]
  [[ "${values[2]}" =~ ^RETURN\ .+$ ]]
}

@test "load available" {

  run plugin_script "load_input.txt"

  echo "$output"

  [ $status -eq $exists ]
}

@test "load unavailable" {  

  if [ "$exists" -ne 0 ]; then
    skip
  fi

  alias git=''

  run plugin_script "load_input.txt"

  [ "$status" -eq 0 ]
}

@test "resolve successfully" {

  if [ "$exists" -ne 0 ]; then
    skip
  fi

  run plugin_resolve "file://${TEST_ROOT_DIRNAME}/data/git"
  
  echo $output

  [ "$status" -eq 0 ]

  resolve_success $output
}

@test "resolve invalid location" {
  if [ "$exists" -ne 0 ]; then
    skip
  fi

  run plugin_resolve "gobbdlygook"

  [ "$status" -ne 0 ]
}

@test "resolve existing" {
  local folder=$(mktemp -u)

  run git clone "${TEST_ROOT_DIRNAME}/data/git" "$folder"

  [ "$status" -eq 0 ]

  run plugin_resolve "${TEST_ROOT_DIRNAME}/data/git" "$folder"

  [ "$status" -eq 0 ]

  resolve_success $output
}

@test "resolve invalid source" {
  skip
}

