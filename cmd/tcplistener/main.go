package main

import (
	"fmt"
	"github.com/brayanMuniz/tcp-to-https/internal/request"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic("悲しい")
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("またね")
		}

		fmt.Println("A new friend has joined, going to print his messages")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("HTTP Version: ", req.RequestLine.HttpVersion)
		fmt.Println("Method: ", req.RequestLine.Method)
		fmt.Println("Target: ", req.RequestLine.RequestTarget)
	}

}
