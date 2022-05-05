#!/bin/bash
mkdir -p data
mkdir -p data/tor
sudo chmod 0700 data/tor -R
sudo chown 100:nogroup data/tor -R
docker-compose build --no-cache