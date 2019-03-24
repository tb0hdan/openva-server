package main

import (
	"github.com/glendc/go-external-ip"
	"github.com/shirou/gopsutil/net"
	gnet "net"
	"strings"
)

func GetExternalIP() (ip gnet.IP, err error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	return consensus.ExternalIP()
}

func GetValidIfaces() (validIfaces []net.InterfaceStat, err error) {
	//
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		// skip ifaces without address
		if len(iface.Addrs) == 0 {
			continue
		}
		// skip ifaces with link-local IPv6 only
		if len(iface.Addrs) == 1 && strings.HasPrefix(iface.Addrs[0].Addr, "fe80::") {
			continue
		}
		// skip ifaces that don't have "up" flag
		up := 0
		for _, flag := range iface.Flags {
			if flag == "up" {
				up++
			}
		}
		if up == 0 {
			continue
		}
		//
		validIfaces = append(validIfaces, iface)
	}
	return
}

func HostPortToIP(hostPort string) (ip gnet.IP, err error) {
	host, _, err := gnet.SplitHostPort(hostPort)
	if err != nil {
		return nil, err
	}
	return gnet.ParseIP(host), nil
}

func ServerIPForClientHostPort(hostPort string) (foundIP string) {
	requestorIP, err := HostPortToIP(hostPort)
	if err != nil {
		panic(err)
	}

	validIfaces, err := GetValidIfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range validIfaces {
		for _, addr := range iface.Addrs {
			parsedAddr, ipnet, err := gnet.ParseCIDR(addr.Addr)
			if err != nil {
				panic(err)
			}
			if ipnet.Contains(requestorIP) {
				foundIP = parsedAddr.String()
				break
			}
		}
		if len(foundIP) > 0 {
			break
		}
	}
	if len(foundIP) == 0 {
		extIP, err := GetExternalIP()
		if err != nil {
			panic(err)
		}
		foundIP = extIP.String()
	}
	return foundIP
}
