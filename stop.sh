#!/bin/bash

docker stop ascii-art-web
docker rm ascii-art-web
docker rmi ascii-art-web-img
docker builder prune