#! /usr/bin/env bash

docker-compose down
docker-compose build
docker-compose -p $EXPOSE_PORT up -d