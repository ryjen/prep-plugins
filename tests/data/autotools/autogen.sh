#!/bin/sh

aclocal \
	&& automake -ac \
	&& autoconf $@