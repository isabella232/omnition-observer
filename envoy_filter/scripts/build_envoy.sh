#!/usr/bin/env bash

set -e

cd /source
# TODO: enabled production optimized builds
bazel build //source:envoy-static

mkdir -p /source/build
cp bazel-bin/source/envoy-static /source/build/envoy

chown $HOST_USERID:$HOST_GROUPID /source/build/envoy
chmod 775 /source/build/envoy 
