#!/usr/bin/env bash
[[ `systemctl status redis | awk '/Active/{print $2}'` == inactive ]] && sudo systemctl start redis
CompileDaemon -build="go build" -recursive="true" -command="./freightsrv"