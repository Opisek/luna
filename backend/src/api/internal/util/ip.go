package util

import (
	"net"
	"regexp"

	"github.com/gin-gonic/gin"
)

var xForwardedForFirstAddress = regexp.MustCompile(`^\s*([^,\s]+)\s*(?:,|$)`)

func DetermineClientAddress(c *gin.Context) net.IP {
	var addr net.IP

	addr = net.ParseIP(c.GetHeader("X-Real-IP"))
	if addr != nil {
		return addr
	}

	addr = net.ParseIP(c.GetHeader("True-Client-IP"))
	if addr != nil {
		return addr
	}

	matches := xForwardedForFirstAddress.FindStringSubmatch(c.GetHeader("X-Forwarded-For"))
	if len(matches) == 2 {
		addr = net.ParseIP(matches[1])
		if addr != nil {
			return addr
		}
	}

	return net.ParseIP(c.ClientIP())
}
