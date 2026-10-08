package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"ping/core"
	"strings"
	"time"

	"golang.org/x/net/ipv4"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the address to check: ")

	rawConn, conn, err := core.SetupSockets()
	if err != nil {
		log.Fatal(err)
		return
	}

	core.EnableFlags(conn, []ipv4.ControlFlags{
		ipv4.FlagTTL,
		ipv4.FlagDst,
		ipv4.FlagInterface,
	})

	defer rawConn.Close()

	addr, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	addr = strings.TrimSpace(addr)

	fmt.Printf("Start of checking %s\n", addr)

	dst := core.ParseIP(addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.Repeat("-", 20))

	for i := 1; i <= 5; i++ {
		msgBytes, msgErr := core.CreateIcmpMessage(0, i, 1)

		if msgErr != nil {
			fmt.Printf("Something went wrong package %d lost \n", i)
			continue
		}

		startAt := time.Now()
		_, sendingErr := rawConn.WriteTo(msgBytes, dst)
		if sendingErr != nil {
			fmt.Printf("Something went wrong package %d cannot be sent \n", i)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		reply := make([]byte, 1500)

		n, cm, peer, readingErr := conn.ReadFrom(reply)
		if readingErr != nil {
			fmt.Printf("Request timeout")
			continue
		}

		duration := time.Since(startAt)

		fmt.Printf("Peer %s, TTL: %d\n", peer, cm.TTL)
		fmt.Printf("Source IP %s, Destination IP %s\n", cm.Src, cm.Dst)
		fmt.Printf("ICMP bytes: %d\n", n)
		fmt.Printf("Duration: %s\n", duration)
		fmt.Println(strings.Repeat("-", 20))

		time.Sleep(2 * time.Second)
	}
}
