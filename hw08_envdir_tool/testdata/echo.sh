#!/usr/bin/env bash

echo -e "HELLO is (${HELLO})
BAR is (${BAR})
FOO is (${FOO})
UNSET is (${UNSET-"_"})
ADDED is (${ADDED})
EMPTY is (${EMPTY})
arguments are $*"
