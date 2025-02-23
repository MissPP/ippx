## ippx proxy

[![Build Status](https://travis-ci.com/MissPP/ippx.svg?branch=main)](https://travis-ci.com/MissPP/ippx)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Introduction

ippx is a sophisticated tunneling proxy designed to facilitate seamless access to restricted networks, enabling you to work from virtually anywhere—whether you're at home or on the go.

# Overview
 
Currently, ippx supports the TCP protocol, with an option to choose SOCKS5 for clients. When utilizing the SOCKS5 protocol, you may configure a connection pool based on TCP, which requires setting up a username and password for access.

Please note that at this time, UDP and IPv6 are not supported.

This project is intended strictly for technical research purposes. Users are advised to employ this tool responsibly and in compliance with applicable laws. The author disclaims any liability for misuse or unlawful actions.
* For technical research purposes only. Please use this project responsibly and in accordance with the law. The author will not be held liable for any consequences arising from misuse.
 
 

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

* Ensure that the ulimit does not impose any restrictions on the calling process.

* If necessary, simply start the servers.

 
## Example
1. OpenSocksProxy:
server socks protocol
```
curl -x socks5://user:pwd@127.0.0.1:7555 yourdomain.com
```
 -- your client (You can consider it as a secondary proxy)
```
curl -x 127.0.0.1:7554 yourdomain.com
 ```
2. CloseSocksProxy: just traffic forwarding
```
curl -x 127.0.0.1:7555 yourdomain.com
```
3. Please refer to ippx/_example for the SOCKS5 request example in Golang (you may need to run go mod tidy).

# License

ippx is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/MissPP/ippx/LICENSE)
