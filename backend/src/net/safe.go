package net

import (
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"
)

// isPublicIP reports whether ip is a globally routable unicast address.
// Loopback, private, link-local (including the 169.254.169.254 cloud metadata
// endpoint), unspecified and multicast addresses are all rejected.
func isPublicIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified() || ip.IsPrivate() || ip.IsMulticast() {
		return false
	}
	return true
}

// guardedControl runs after DNS resolution on the concrete IP that is about to
// be dialed. Because it inspects the resolved address, it defends against SSRF
// via internal hostnames, DNS rebinding and redirect chains to private targets.
func guardedControl(network, address string, _ syscall.RawConn) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("could not parse dial address %q", address)
	}
	if !isPublicIP(ip) {
		return fmt.Errorf("refusing to connect to non-public address %v", ip)
	}
	return nil
}

var guardedDialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
	Control:   guardedControl,
}

// GuardedClient is an http.Client that refuses to connect to loopback, private,
// link-local or otherwise non-public addresses. It must be used for every
// outbound request whose destination is derived from user input (remote
// calendar sources) to prevent server-side request forgery.
var GuardedClient = &http.Client{
	Transport: &http.Transport{
		DialContext:           guardedDialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}
