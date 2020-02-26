package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	hostPort := os.Args[1]
	l, err := net.Listen("tcp", hostPort)
	if err != nil {
		panic(err)
	}

	for {
		log.Println("Listening", hostPort)
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		log.Println("Accepted")
		scanner := bufio.NewScanner(conn)
		ok := scanner.Scan()
		if ok {
			s := scanner.Text()
			r := strings.ToUpper(s) + "\n"
			log.Print(r)
			conn.Write([]byte(r))
			conn.Close()
			if s == "*" {
				log.Println("Exiting")
				break
			}
		}
	}
}
