package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("host:port must be specified")
		return
	}
	addr := os.Args[1]
	fmt.Printf("TCP client dailing to %s ...\n", addr)
	c, err := net.Dial("tcp", addr)
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
