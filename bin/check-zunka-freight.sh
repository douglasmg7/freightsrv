#!/usr/bin/env bash
read -r HOST USER PASS <<< $(./auth.sh zunkasite)

curl -u $USER:$PASS -X GET $HOST/freightsrv/freights/zunka \
    -H "Content-Type: application/json" \
	-d '{
		"CEPDestiny": "31170210",
		"Weight":     1500,
		"Length":     20,
		"Height":     30,
		"Width":      40,
        "Price":      90.54
	}'
printf "\n"
