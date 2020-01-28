#!/usr/bin/env bash
[[ `systemctl status redis | awk '/Active/{print $2}'` == inactive ]] && sudo systemctl start redis
CompileDaemon -build="go build" -recursive="true" -command="./freightsrv dev"

# CompileDaemon -build="go build" -include="*.tpl" -include="*.tmpl" -include="*.gohtml" -include="*.css" -recursive="true" -command="./zunkasrv dev"
# go run *.go dev