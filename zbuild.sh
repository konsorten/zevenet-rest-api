#!/bin/sh

set -e

export GOARCH=386

if [ -f ./zevenet-rest-api ]; then
    rm -f ./zevenet-rest-api
fi

~/go/bin/ktn-build-info

~/go/bin/swag init --swagger ./www/swagger/ --generalInfo version.go

rm -f -R ./docs

go build .

if [ ! -x ./zevenet-rest-api ]; then
    echo "BUILD FAILED"
    exit 1
fi

# copy files
sudo cp -f ./zevenet-rest-api /usr/local/zevenet/bin/zevenet-rest-api
sudo cp -f -R ./cherokee/* /usr/local/zevenet/app/cherokee/
sudo cp -f -R ./www/* /usr/local/zevenet/www/
