#!/usr/bin/env bash

set -e
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

mkdir -p config/wg_confs/
cp  config-files/server.wg0.conf config/wg_confs/wg0.conf
#chown -R 911:911 config/

docker compose up -d --build
docker compose logs -f
