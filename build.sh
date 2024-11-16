#!/bin/bash

docker build -t ascii-art-web-img .
docker create --name ascii-art-web -p 8080:8080 ascii-art-web-img
docker images
docker ps -a
