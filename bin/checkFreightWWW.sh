#!/usr/bin/env bash
curl -u test:1234 -X GET https://www.zunka.com.br/freightsrv/freights/zunka \
    -H "Content-Type: application/json" \
	-d '{
		"CEPDestiny": "31170210",
		"Weight":     1500,
		"Length":     20,
		"Height":     30,
		"Width":      40
	}'
