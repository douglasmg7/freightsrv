#!/usr/bin/env bash
# clear
[[ `systemctl status redis | awk '/Active/{print $2}'` == inactive ]] && sudo systemctl start redis
# go test
# go test -run SaveFreightRegion

# go test -run Handler -args -dev=true

# Clean cache.
if [[ $1 == "--clean-cache" ]]; then
    echo Cleaning cache...
    KEYS=`redis-cli keys freightsrv-*`
    [[ ! -z $KEYS ]] && redis-cli del $KEYS
    exit
    # echo $KEYS
    # redis-cli del `redis-cli keys freightsrv-*`
fi

# Reset db.
if [[ $1 == "--clean-db" ]]; then
    echo Cleaning db...
    cd bin/db
    ./rcp.sh
    cd ../..
    exit
fi

# go test -run otoboy
go test