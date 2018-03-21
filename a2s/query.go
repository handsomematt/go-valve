package a2s

import (
	"log"
	"net"
	"strconv"
	"time"
)

type Querier struct {
	udpConnection *net.UDPConn
}

func NewQuerier(host string, port int, timeout time.Duration) *Querier {
	service := host + ":" + strconv.Itoa(port)
	remoteAddr, err := net.ResolveUDPAddr("udp", service)

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Timeout
	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		log.Fatal(err)
	}

	return &Querier{udpConnection: conn}
}

// Close ...
func (querier *Querier) Close() {
	querier.udpConnection.Close()
}
