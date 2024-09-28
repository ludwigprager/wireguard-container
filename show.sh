#!/usr/bin/env bash

set -e
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

docker exec -ti wireguard wg show all dump
