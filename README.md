# portfwd
A tiny, configurable port forwarder written in Go

## Install
`go get github.com/machinae/portfwd/...`

## Configure
Put forwardings you would like to configure in a TOML config file.
portfwd.cfg
```toml
# Forwardings can have any name you want
# Make sure to quote strings
[development]
from = ":8000"
to = "remote:8000"

[staging]
from = ":9000"
# Forwarding hosts are rotated randomly by default
to = [
"remote1.com:9001",
"remote2.com:9001"
]

[production]
# Set strategy to roundrobin to rotate forwarders in order
strategy = "roundrobin"
from = ":8000"
to = [
"1.2.3.4:8000",
"5.6.7.8:8000"
]

## Run
`portfwd -c portfwd.cfg`
```
