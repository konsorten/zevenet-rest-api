#!/bin/sh

set -e

export GOARCH=386

if [ -f ./zevenet-rest-api ]; then
    rm -f ./zevenet-rest-api
fi

go build .

if [ ! -x ./zevenet-rest-api ]; then
    echo "BUILD FAILED"
    exit 1
fi

# copy files
sudo cp ./zevenet-rest-api /usr/local/zevenet/bin/zevenet-rest-api
sudo cp ./cherokee_rest_api.conf /usr/local/zevenet/app/cherokee/etc/cherokee/cherokee_rest_api.conf
