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

	c, err := net.Dial("tcp", os.Args[1])
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
