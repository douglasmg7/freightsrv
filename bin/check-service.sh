#!/usr/bin/env bash
read -r HOST USER PASS <<< $(./auth.sh zunka)

curl -u $USER:$PASS $HOST/freightsrv/hello
printf "\n"