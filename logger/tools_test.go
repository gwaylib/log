package logger

import (
	"testing"
)

func TestGetIp(t *testing.T) {
	ip := getIp()
	println(ip)
}

func TestGetServerName(t *testing.T) {
	println(getServerName())
}
