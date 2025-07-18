package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	//  go run ./cmd/tcplistner | tee /tmp/tcp.txt
	//  printf "好き" | nc -c -w 1 127.0.0.1 42069

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
		channelString := getLinesChannel(conn)
		for line := range channelString {
			fmt.Print(line)
		}

	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineChannel := make(chan string)

	go func() {
		line := ""
		buffer := make([]byte, 8)
		for {
			amountRead, err := f.Read(buffer)
			if err == io.EOF {
				if line != "" {
					lineChannel <- line
				}
				close(lineChannel)
				f.Close()
				return
			}

			str := string(buffer[:amountRead])
			newline := strings.Index(str, "\n")
			if newline != -1 {
				line += str[0:newline]
				lineChannel <- line

				line = str[newline:len(str)]
			} else {
				line += str
			}
		}
	}()

	return lineChannel
}
