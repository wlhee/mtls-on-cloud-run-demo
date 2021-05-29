package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 1024)
	for {
		_, err := bufio.NewReader(c).Read(buf)
		if err != nil {
			break
		}
		fmt.Print(string(buf))
	}
}
