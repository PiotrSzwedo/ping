package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the address to check: ")
	addr, err := reader.ReadString('\n')

	if err != nil {
		return
	}

	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Printf("Error: %v", closeErr)
		}
	}()

	fmt.Printf("Start of checking %s\n", addr)

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid(),
			Seq:  1,
			Data: []byte("HELLO-GO"),
		},
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	startAt := time.Now()

	ip := net.ParseIP(addr)
	dst := &net.IPAddr{IP: ip}

	//if err != nil {
	//	log.Fatal(err)
	//}

	if _, err := conn.WriteTo(msgBytes, dst); err != nil {
		log.Fatal(err)
	}

	reply := make([]byte, 1500)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		log.Fatalf("Timeout: %v", err)
	}
	duration := time.Since(startAt)

	fmt.Printf("Response %d bytes from %s after %s seconds \n", n, peer, duration)
}
