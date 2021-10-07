#!/usr/bin/bash

if [ "$1" == 'build' ]; then
  docker build -t twi-meteor .
elif [ "$1" == 'run' ]; then
  docker run -d  --env-file .env twi-meteor
elif [ "$1" == 'log' ]; then
  docker exec "$2" tail -f /var/log/app.log
else
  echo 'usage: bash script.sh <build|run|log> <if log, then hash>'
fi