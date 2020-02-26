package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}
	log.Println("connecting ", tcpAddr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	s := strings.Join(os.Args[2:], " ") + "\n"
	_, err = conn.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(reply[0:n]))
}
