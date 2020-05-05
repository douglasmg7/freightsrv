#!/usr/bin/env bash
read -r HOST USER PASS <<< $(./auth.sh zoombuscape)

# Run mode.
if [[ $RUN_MODE == production ]];then
    SKU1=5e60eed63d13910785412eba
    SKU2=5bcb336a4253f81781faca09
else
    SKU1=5e60eed63d13910785412eba
    SKU2=5bcb336a4253f81781faca09
fi

generate_data()
{
  cat <<EOF
{
    "zipcode": "31170210",
    "items": [
        { 
            "amount": 1,
            "sku": "1I3eF1",
            "price": 34.4, 
            "weight": 2,
            "height": 0.3,
            "width": 0.2,
            "length": 0.4
        },
        {
            "amount": 2,
            "sku": "2R343",
            "price": 34.4, 
            "weight": 2,
            "height": 0.3,
            "width": 0.2,
            "length": 0.4
        }   
    ]
}
EOF
}

# curl -u $USER:$PASS -X POST $HOST/freightsrv/freights/zoom \
    # -H "Content-Type: application/json" \
    # -d "$(generate_data)"

CMD="curl -v -u $USER:$PASS $HOST/freightsrv/freights/zoom -H \"Content-Type: application/json\" -d '$(generate_data)'"

# printf "%s\n\n" "$CMD"
echo $CMD
echo
eval $CMD
printf "\n"
