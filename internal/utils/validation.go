package utils

import (
	"net"
	"net/url"
)

func IsValidIPAddress(s string) bool {
	return net.ParseIP(s) != nil
}

func IsValidURL(s string) bool {
	if _, err := url.ParseRequestURI(s); err != nil {
		return false
	}
	return true
}
