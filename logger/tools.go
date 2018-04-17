package logger

import (
	"os"
)

var HostIp = getHost()

func getHost() string {
	hostName, _ := os.Hostname()
	return hostName
}
