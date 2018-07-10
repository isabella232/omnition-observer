#!/bin/bash

set -e

umask 022

if ! getent passwd omnition-proxy >/dev/null; then
    groupadd --gid 1337 omnition-proxy
    useradd --uid 1337 --gid 1337 -d /var/lib/omnition omnition-proxy
    echo "omnition-proxy ALL=NOPASSWD: ALL" >> /etc/sudoers
fi

mkdir -p /var/lib/omnition/envoy
mkdir -p /var/lib/omnition/proxy
mkdir -p /var/lib/omnition/config
mkdir -p /var/log/omnition

chown omnition-proxy.omnition-proxy /var/lib/omnition/envoy /var/lib/omnition/config /var/log/omnition /var/lib/omnition/proxy
echo "almost there"
chmod o+rx /usr/local/bin/envoy

# envoy may run with effective uid 0 in order to run envoy with
# CAP_NET_ADMIN, so any iptables rule matching on "-m owner --uid-owner
# omnition-proxy" will not match connections from those processes anymore.
# Instead, rely on the process's effective gid being omnition-proxy and create a
# "-m owner --gid-owner omnition-proxy" iptables rule in prepare_proxy.sh.
chmod 2755 /usr/local/bin/envoy
chgrp omnition-proxy /usr/local/bin/envoy

echo "starting envoy"
sg omnition-proxy -c "envoy -c /etc/envoy.yaml --v2-config-only"

