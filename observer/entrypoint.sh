#!/bin/bash

export PATH=$PATH:/usr/local/go/bin

USER_ID=${USER_ID:-1001}

useradd --shell /bin/bash -u $USER_ID -o -c "" -m user
export HOME=/home/user

exec /usr/local/bin/gosu user "$@"
