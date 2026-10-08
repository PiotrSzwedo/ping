package core

import (
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func SetupSockets() (net.PacketConn, *ipv4.PacketConn, error) {
	rawConn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")

	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	conn := ipv4.NewPacketConn(rawConn)

	return rawConn, conn, nil
}

func ParseIP(address string) *net.IPAddr {
	ip := net.ParseIP(address)
	dst := &net.IPAddr{IP: ip}
	return dst
}

func CreateIcmpMessage(code int, seq int, dataSize int) ([]byte, error) {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: code,
		Body: &icmp.Echo{
			ID:   os.Getpid(),
			Seq:  seq,
			Data: make([]byte, dataSize),
		},
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return nil, err
	}

	return msgBytes, nil
}

func EnableFlags(conn *ipv4.PacketConn, cfs []ipv4.ControlFlags) {
	manageFlags(conn, cfs, true)
}

func DisableFlags(conn *ipv4.PacketConn, cfs []ipv4.ControlFlags) {
	manageFlags(conn, cfs, false)
}

func manageFlags(conn *ipv4.PacketConn, cfs []ipv4.ControlFlags, action bool) {
	for _, cf := range cfs {
		err := conn.SetControlMessage(cf, action)
		if err != nil {
			log.Fatal(err)
		}
	}
}
