package utils

import (
	"net"
	"net/url"
)

func IsValidIPAddress(s string) bool {
	return net.ParseIP(s) != nil
}

func IsValidURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	_ = u
	if err != nil {
		return false
	}
	return true
}
