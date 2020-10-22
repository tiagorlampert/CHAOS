package models

import "net"

type Device struct {
	Connection     net.Conn `json:"-"`
	Hostname       string   `json:"hostname"`
	Username       string   `json:"username"`
	UserID         string   `json:"user_id"`
	OSName         string   `json:"os_name"`
	MacAddress     string   `json:"mac_address"`
	LocalIPAddress string   `json:"local_ip_address"`
	Port           string   `json:"port"`
	FetchedUnix    int64    `json:"fetched_unix"`
}
