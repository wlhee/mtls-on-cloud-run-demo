package main

import (
	"fmt"
	"net"
)

func main() {
	addr := ":7777"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	fmt.Printf("Demo TCP server starts listening on %s\n", addr)

	for {
		conn, err := l.Accept()
		fmt.Println("new connection accepted.")
		if err != nil {
			fmt.Printf("connection error: %v\n", err)
			return
		}

		go func() {
			conn.Write([]byte("================================================================\n"))
			conn.Write([]byte("== Congrats!                                                  ==\n"))
			conn.Write([]byte("== If you see this message, it means                          ==\n"))
			conn.Write([]byte("== you've successfully run the mTLS demo on Google Cloud Run! ==\n"))
			conn.Write([]byte("================================================================\n"))
			conn.Close()
		}()
	}
}
