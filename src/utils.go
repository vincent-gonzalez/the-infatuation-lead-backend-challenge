package main

import (
	"fmt"
	"net"
)

func CreateConnection(protocol string, url string, port uint) (net.Conn, error) {
	newConnection, err := net.Dial(protocol, fmt.Sprintf("%s:%d", url, port))
	if err != nil {
		return nil, err
	}
	return newConnection, err
}
