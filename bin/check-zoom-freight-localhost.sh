#!/usr/bin/env bash
curl -u test:1234 -X GET localhost:8081/freightsrv/freights/zoom \
    -H "Content-Type: application/json" \
	-d '{
		"CEPDestiny": "31170210",
		"productId":  "5e60eed63d13910785412eba"
	}'
printf "\n"

# "productId":  "5c19eab2fbed5f0a1c19dcc8"
# "productId":  "5e60eed63d13910785412eba"  # Aldo product.