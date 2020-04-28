#!/usr/bin/env bash

# Run mode.
[[ $RUN_MODE == production ]] && HOST="https://www.zunka.com.br" || HOST="localhost:8081"

case $1 in
  zoom)
	USER=zoom-buscape
	PASS=rt54J9_W29
    ;;

  zunka)
    USER=zunka
	PASS=fi4vI9_qAA
    ;;

  browser)
    USER=zunka-browser
	PASS=3j-8clHd_4
    ;;

  *)
    USER=user
	PASS=pass
    ;;
esac

echo $HOST $USER $PASS
