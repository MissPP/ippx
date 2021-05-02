## ippx proxy

[![Build Status](https://travis-ci.com/MissPP/ippx.svg?branch=main)](https://travis-ci.com/MissPP/ippx)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Introduction

ippx is a tunnel proxy which can help you get through intranet. So you can work at home now !

# Overview
 
 ippx now support tcp protocol ,also you can choose the socks5 protocol and start your client with a connection pool based on tcp
 had to set the username and password when you use socks protocol 
 
 doesn't support udp and ip6 for now  
* For technical research only, please use this project reasonably and legally, otherwise the author will not be responsible for the consequences
 
 

# Install

install golang 

go env -w GO111MODULE=auto


## Encrypt
RC4 (change it by yourself if you need to)

# Setting
config path
```
  ▾ ippx/
    ▾ config/
        config.json
```

set your server and client address in the file


# Getting Started
**In your server**
```
  ▾ ippx/
    ▾ cmd_server/
```

    go run start_server.go

**In your client**
```
  ▾ ippx/
    ▾ cmd_client/
```

    go run start_client.go

* Make sure the ulimit will not  set some limit for the calling process

* If you need to, just start the servers

 
## Example
1. OpenSocksProxy:
server socks protocol
```
curl -x socks5://user:pwd@127.0.0.1:7555 yourdomain.com
```
 -- your client (You can think of it as a secondary proxy)
```
curl -x 127.0.0.1:7554 yourdomain.com
 ```
2. CloseSocksProxy: just traffic forwarding
```
curl -x 127.0.0.1:7555 yourdomain.com
```
3. Please see the `ippx/_example`for the socks5 request by golang ( you may need `go mod tidy` )

# License

ippx is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/MissPP/ippx/LICENSE)
