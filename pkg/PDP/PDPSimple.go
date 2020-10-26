package PDPSimple

import (
	"context"
	"fmt"
	"log"
	"net"
)

func listenMultiCast(ctx context.Context, address string) {

	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(err)
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	maxDatagramSize := 65535

	conn.SetReadBuffer(maxDatagramSize)

	// Loop forever reading from the socket
	for {
		select {
		case <-ctx.Done():
			log.Printf("multicast listening exit...")
		default:
			buffer := make([]byte, maxDatagramSize)
			numBytes, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Fatal("ReadFromUDP failed:", err)
			}

			fmt.Printf("<%s> %s \n", remoteAddr, buffer[:numBytes])
		}
	}
}
