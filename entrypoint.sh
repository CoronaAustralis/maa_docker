#!/bin/bash

if [ "$(ls -A /app/config)" = "" ]; then
    echo "empty folder, copy config"
    cp -r  /tmp/config/* /app/config/
fi

if [ $PROXY ]; then
    git config --global http.proxy $PROXY
    git config --global https.proxy $PROXY
    export HTTP_PROXY=$PROXY
    export HTTPS_PROXY=$PROXY
fi

/app/main