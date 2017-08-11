#!/bin/bash
# waiter.sh

set -e

host="$1"
shift
cmd="$@"

sleep 5

exec $cmd
