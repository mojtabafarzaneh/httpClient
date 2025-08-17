package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mojtabafarzaneh/httpClinet/request"
)

func main() {
	listenter, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}

	for {
		conn, err := listenter.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)

		}
		fmt.Printf("Request Line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	}
}
