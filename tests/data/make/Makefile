
FILE=make_output.log
SHELL=bash

all: install

install: run test

run:
	@echo -n "Making test file... "
	@touch ${FILE}

test:
	@test -e make_output.log && echo "Succeeded" || exit 1
