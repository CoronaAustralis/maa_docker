#!/bin/bash

if [ "$(ls -A /app/)" = "" ]; then
    echo "empty folder, copy config"
    cp -r /source/* /app/
fi

if [ $http_proxy ]; then
    git config --global http.proxy $http_proxy
    echo "add git http.proxy"
fi

if [ $https_proxy ]; then
    git config --global https.proxy $https_proxy
    echo "add git https.proxy"
fi

if [ ! -f "/root/.local/share/maa/lib/libMaaCore.so" ]; then
    maa install
fi

supercronic /app/crontab