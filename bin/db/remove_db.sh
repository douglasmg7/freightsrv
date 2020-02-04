#!/usr/bin/env bash 

# ZUNKAPATH not defined.
if [ -z "$ZUNKAPATH" ]; then
	printf "error: ZUNKAPATH not defined.\n" >&2
	exit 1 
fi

# ZUNKA_FREIGHT_DB not defined.
if [ -z "$ZUNKA_FREIGHT_DB" ]; then
	printf "error: ZUNKA_FREIGHT_DB not defined.\n" >&2
	exit 1 
fi

DB=$ZUNKAPATH/db/$ZUNKA_FREIGHT_DB
# DB not exist.
if [[ ! -f $DB ]]; then
    printf "DB %s not exist\n" $DB
    exit 0
fi

printf "Removing db %s\n" $DB
rm $DB
