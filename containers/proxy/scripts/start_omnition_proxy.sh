#!/bin/bash

set -e

umask 022

export OBS_TRACING_DRIVER=$TRACING_DRIVER
export OBS_TRACING_ADDRESS=$TRACING_ADDRESS
export OBS_TRACING_PORT=$TRACING_PORT

export OBS_TLS_ENABLED=$TLS_ENABLED
export OBS_TLS_CERT=$TLS_CERT
export OBS_TLS_KEY=$TLS_KEY
export OBS_TLS_CA_CERT=$TLS_CA_CERT

export OBS_ADMIN_PORT=$ADMIN_PORT
export OBS_ADMIN_LOG_PATH=$ADMIN_LOG_PATH

export OBS_INGRESS_PORT=$INGRESS_PORT
export OBS_EGRESS_PORT=$EGRESS_PORT

export SERVICE_NAME=${SERVICE_NAME:-'unknown-service'}


TODO(owais): Test system wide CA cert approval 
#if [ -z "$OBS_CA_CERT" ]; then
#  echo $OBS_CA_CERT > /usr/local/share/ca-certificates/envoy.crt
#  update-ca-certificates
#fi

echo "setting up roles"
if ! getent passwd omnition-proxy >/dev/null; then
    groupadd --gid 1337 omnition-proxy
    useradd --uid 1337 --gid 1337 -d /var/lib/omnition omnition-proxy
    echo "omnition-proxy ALL=NOPASSWD: ALL" >> /etc/sudoers
fi

mkdir -p /var/lib/omnition/envoy
mkdir -p /var/lib/omnition/proxy
mkdir -p /var/lib/omnition/config
mkdir -p /var/lib/omnition/tls
mkdir -p /var/log/omnition/

echo "setting up permissions"
chown -R omnition-proxy:omnition-proxy /var/lib/omnition/ /var/log/omnition
chmod o+rx /usr/local/bin/envoy

# envoy may run with effective uid 0 in order to run envoy with
# CAP_NET_ADMIN, so any iptables rule matching on "-m owner --uid-owner
# omnition-proxy" will not match connections from those processes anymore.
# Instead, rely on the process's effective gid being omnition-proxy and create a
# "-m owner --gid-owner omnition-proxy" iptables rule in prepare_proxy.sh.
chmod 2755 /usr/local/bin/envoy
chgrp omnition-proxy /usr/local/bin/envoy

observer > /etc/envoy.yaml

if [ $1 = "show-config" ];
  then
  cat /etc/envoy.yaml
elif [ $1 = "run" ]
  then
  echo "starting envoy"
  sg omnition-proxy -c "envoy -c /etc/envoy.yaml -l info --v2-config-only --service-cluster $SERVICE_NAME"
else
  $1
fi
