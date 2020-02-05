#!/usr/bin/env bash
clear
[[ `systemctl status redis | awk '/Active/{print $2}'` == inactive ]] && sudo systemctl start redis
# go test
go test -run FreightRegion