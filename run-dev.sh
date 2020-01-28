#!/usr/bin/env bash
[[ `systemctl status redis | awk '/Active/{print $2}'` == inactive ]] && sudo systemctl start redis
go build && ./freightsrv