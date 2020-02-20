#!/usr/bin/env bash
# clear
[[ `systemctl status redis | awk '/Active/{print $2}'` == inactive ]] && sudo systemctl start redis
# go test
# go test -run SaveFreightRegion

# go test -run Handler -args -dev=true

# Default clean cache.
if [[ $1 != "--cache" ]]; then
    echo Cleaning cache...
    KEYS=`redis-cli keys freightsrv-*`
    [[ ! -z $KEYS ]] && redis-cli del $KEYS
    # echo $KEYS
    # redis-cli del `redis-cli keys freightsrv-*`
fi

go test -run UpdateMotoboyFreight