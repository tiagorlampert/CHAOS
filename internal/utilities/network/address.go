package network

import (
	"log"
	"net"
)

func GetLocalIP() net.IP {
	conn, err := net.Dial(`udp`, `8.8.8.8:80`)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
