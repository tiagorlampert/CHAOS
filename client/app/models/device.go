package models

type Device struct {
	Hostname       string `json:"hostname"`
	Username       string `json:"username"`
	UserID         string `json:"user_id"`
	OSName         string `json:"os_name"`
	MacAddress     string `json:"mac_address"`
	LocalIPAddress string `json:"local_ip_address"`
	FetchedUnix    int64  `json:"fetched_unix"`
}
