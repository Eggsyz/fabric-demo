#!/bin/bash
if [ "$1" == "gm" ]; then
echo "start fabric by bccsp: $1"
docker-compose -f docker-compose-gm.yaml up -d
else
echo "start fabric by default bccsp: sw"
docker-compose -f docker-compose-sw.yaml up -d
fi
