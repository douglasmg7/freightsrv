#!/usr/bin/env bash
read -r HOST USER PASS <<< $(./auth.sh zunkasite)

curl -u $USER:$PASS $HOST/freightsrv/hello
printf "\n"