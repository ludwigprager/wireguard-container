#!/usr/bin/env bash

set -eu

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $BASEDIR


CONF=wg.yaml

if [[ -f $CONF ]]; then
  echo "wireguard config file already exists"
else
  echo "Creating wireguard config file"
  cp wg.yaml-vm $CONF
fi


sed -i "s|SERVER_IP|$(curl ipinfo.io/ip)|" $CONF
sed -i "s|SERVER_KEY|$(wg genkey)|" $CONF
sed -i "s|CLIENT1_KEY|$(wg genkey)|" $CONF
sed -i "s|CLIENT2_KEY|$(wg genkey)|" $CONF
