# TheStoneProxy
Lightweight go-written L4 proxy for minecraft bedrock editon (mcpe/mcbe).

## Features

- Rest-API to get players info.
- Prometheus server to collect metrics.
- Flexible configuration.
- Using the Raknet protocol, so it doesn't need to be updated every time new mcpe version released.
- Ready to rock Docker.

## Getting started

### Build on your own

*You must have **go** installed on your PC*

Clone git repo and get into project dir
```bash
git clone https://github.com/OpenFarLands/TheStoneProxy && cd TheStoneProxy
```

And build it
```bash
go build .
```

Run proxy
```bash
./TheStoneProxy
```

### Use Docker

*You must have **docker** installed on your PC*

Download default config 
```bash
wget https://raw.githubusercontent.com/OpenFarLands/TheStoneProxy/main/example.config.toml

```

Start the container
```bash
docker run -v ${PWD}/example.config.toml:/app/config.toml --network host --restart unless-stopped mikhaillav/stone-proxy:latest

```

*Recommending to use **host** network type*

### Use Docker-Compose

*You must have **docker** & **docker-compose** installed on your PC*

Example of *docker-compose.yml* for 2 proxy-servers
```yaml
version: '3.7'

services:
    proxy1:
        image: mikhaillav/stone-proxy:latest
        volumes:
            - ${PWD}/first.config.toml:/app/config.toml
        restart: unless-stopped
        network_mode: host


    proxy2:
        image: mikhaillav/stone-proxy:latest
        volumes:
            - ${PWD}/second.config.toml:/app/config.toml
        restart: unless-stopped
        network_mode: host

```

*Recommending to use **host** network type*


After creating docker-compose.yml and config files start your proxy
```bash
docker-compose up -d
```

## Configuration

Example of config

```toml
[Network]
# Address of proxy for clients to connect
LocalAddress = "0.0.0.0:19132"
# Address of origin server for proxy to relay
RemoteAddress = "localhost:19133"
# Time in seconds of inactivity before connection will be closed
Timeout = 60
# Time in second of waiting before sending ping packet (and get motd) to origin server
MotdGetInterval = 10

[OfflineMotd]
# Motd that proxy will retrun if upstream server if online (or send invalid pong packet)
Motd = '§c§lOffline'
# Protocol version (e.g. 618) that proxy will return if upstream server if online (or send invalid pong packet)
ProtocolVersion = 1
# Version name (e.g. 1.20.31) that proxy will return if upstream server if online (or send invalid pong packet)
VersionName = '1.0.0'
# Level name that proxy will return if upstream server if online (or send invalid pong packet)
LevelName = 'Powered by TheStoneProxy'

[Api]
# Whether use api server
UseApiServer = false
# Address of api server
ApiServerAddress = ":3030"
# Array of ip addresses that can use api. Keep it empty to allow any ip address
ApiWhitelist = ["[::1]"]

[Metrics]
# Whether use prometheus server
UsePrometheus = false
# Address of prometheus server
PrometheusAddress = ":3031"
# Token for auth request to metrics server. Keep it empty to disable auth
PrometheusBearerAuthToken = "helloworld"
```

## Licenese

This project is doing under MIT License!
