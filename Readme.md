# PasteTor

A simple Paste Tor service with a small management overlay written in golang + a redis database for high performance
needs.
Allows quick setup for your own hidden service.

### Requirements

Docker + Docker-Compose

Around 200mb space for the docker images. After building the huge buster image can be removed.

### Setup

```bash
sudo bash rebuild.sh
```

## Start

```bash
sudo docker-compose up -d
# Wait 10 seconds and get your hostname:
sudo cat data/tor/hidden_service/hostname
```