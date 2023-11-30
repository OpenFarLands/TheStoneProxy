# TheStoneProxy
Lightweight go-written L4 proxy for minecraft bedrock editon (mcpe/mcbe).

## Features

- Rest-API to get players info.
- Prometheus server to collect metrics.
- Flexible configuration.
- Using the Raknet protocol, so it doesn't need to be updated every time new mcpe version released.
- Ready to rock Dockerfile.

## Getting started

### Build by your own

*You must have go installed on your PC*

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

*Comming soon...*

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
# Unconnected pong message that will return if proxy failed to fetch that from origin server
# SERVER_UNIQUE_ID, PORT_V_4, PORT_V_6 WILL BE REPLACED TO CORRECT ONES
OfflinePongMessage = "MCPE;§c§lOffline !(;618;1.20.31;0;10;SERVER_UNIQUE_ID;Bedrock level?;Survival;1;PORT_V_4;PORT_V_6;1;"

[Api]
# Whether use api server
UseApiServer = true
# Address of api server
ApiServerAddress = ":3030"
# Array of ip addresses that can use api. Keep it empty to allow any ip address
ApiWhitelist = ["[::1]"]

[Metrics]
# Whether use prometheus server
UsePrometheus = true
# Address of prometheus server
PrometheusAddress = ":3031"
# Token for auth request to metrics server. Keep it empty to disable auth
PrometheusBearerAuthToken = "helloworld"
````

## Licenese

This project is doing under MIT License!