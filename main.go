package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./m.txt")
	if err != nil {
		log.Fatal("Could not read file")
	}

	channelString := getLinesChannel(file)
	for line := range channelString {
		fmt.Print(line)
	}

	fmt.Println("---")
	fmt.Println("バイバイ")
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
