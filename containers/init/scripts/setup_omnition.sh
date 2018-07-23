#!/bin/bash

set -e

umask 022

if ! getent passwd omnition-proxy >/dev/null; then
    groupadd --gid 1337 omnition-proxy
    useradd --uid 1337 --gid 1337 -d /var/lib/omnition omnition-proxy
    echo "omnition-proxy ALL=NOPASSWD: ALL" >> /etc/sudoers
fi

echo "Removing old iptables rules"
# Remove the old chains, to generate new configs.
iptables -t nat -D PREROUTING -p tcp -j OMNITION_INBOUND 2>/dev/null || true
iptables -t nat -D OUTPUT -p tcp -j OMNITION_OUTPUT 2>/dev/null || true

# Flush and delete the omnition chains.
iptables -t nat -F OMNITION_OUTPUT 2>/dev/null || true
iptables -t nat -X OMNITION_OUTPUT 2>/dev/null || true
iptables -t nat -F OMNITION_INBOUND 2>/dev/null || true
iptables -t nat -X OMNITION_INBOUND 2>/dev/null || true
# Must be last, the others refer to it
iptables -t nat -F OMNITION_REDIRECT 2>/dev/null || true
iptables -t nat -X OMNITION_REDIRECT 2>/dev/null || true


echo "Add new iptables rules"
# new chain to envoy port
iptables -t nat -N OMNITION_REDIRECT
iptables -t nat -A OMNITION_REDIRECT -p tcp -j REDIRECT --to-port 15001

# Inbound traffic
iptables -t nat -N OMNITION_INBOUND
iptables -t nat -A PREROUTING -p tcp -j OMNITION_INBOUND
iptables -t nat -A OMNITION_INBOUND -p tcp -j OMNITION_REDIRECT

# Outbound traffic
iptables -t nat -N OMNITION_OUTPUT
iptables -t nat -A OUTPUT -p tcp -j OMNITION_OUTPUT

# Ignore outbound traffic originating from envoy process' gid
# iptables -t nat -A OMNITION_OUTPUT -m owner --gid-owner omnition-proxy -j RETURN
iptables -t nat -A OMNITION_OUTPUT -m owner --gid-owner omnition-proxy -j RETURN

# redirect remaining outbound packets
iptables -t nat -A OMNITION_OUTPUT -j OMNITION_REDIRECT


echo "Omnition init complete"
