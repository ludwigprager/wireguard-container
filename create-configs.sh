#!/usr/bin/env bash

set -eu

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $BASEDIR

function go-in-docker() {
  local command="$*"

  docker run -ti --rm \
    -w /work \
    -v $(pwd):/work/ \
    -e GOMODCACHE=/work/go/ \
    -e GOCACHE=/work/go/build-cache \
    golang:1.23.1 \
    $command
}

if [[ ! -f go.mod ]]; then
  go-in-docker go mod init wg-config
fi

if [[ ! -f go.sum ]]; then
  go-in-docker go get \
    gopkg.in/yaml.v2 \
    github.com/davecgh/go-spew/spew \
    golang.zx2c4.com/wireguard/wgctrl/wgtypes
fi

go-in-docker go build main.go


mkdir -p config-files
go-in-docker ./main config-files



# cp server.wg0.conf wg0.conf

for i in config-files/*.conf; do
  docker run --rm -t \
    --user="$(id -u):$(id -g)" --net=none \
    -v "$(pwd):/tmp" leplusorg/qrcode \
    qrencode -l L -r /tmp/$i -o /tmp/$i.png

done
