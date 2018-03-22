package logger

import (
	"net"
	"strings"
)

var HostIp = getIp()
var ServerName = getServerName()

func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, addr := range addrs {
		ip := strings.Split(addr.String(), "/")[0]
		code := strings.Split(ip, ".")
		switch code[0] {
		case "127":
			continue
		default:
			return ip
		}
	}
	return ""
}

func getServerName() string {
	// return "etlog"
	return getIp()
}
