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
	printf "error: DB %s not exist.\n" $DB >&2
    exit 1
fi

echo "Inserting data..."
sqlite3 $DB < $(dirname $0)/data.sql
