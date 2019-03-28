#!/bin/bash

# https://raw.githubusercontent.com/grandcat/zeroconf/master/examples/register/server.go
go run server.go -service _openva._tcp -name OpenVA-01 -port 50001
