package multicast

import (
	"fmt"
	"net"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const _UDP4 = "udp4"


type UdpBroadcast struct {

	conn *net.UDPConn;
}

func (sock *UdpBroadcast) Send(bytes []byte) {

	fmt.Println("Sending ", bytes)
	sent, err := sock.conn.Write(bytes)
	if err != nil {
		fmt.Println("Failed to send bytes")
	}

	fmt.Printf("Sent %d Bytes \n", sent)
}

func (sock *UdpBroadcast) SendProto(msg protoreflect.ProtoMessage) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	sock.Send(bytes)
}

func NewBroadcast(address string) *UdpBroadcast{
	addr, err := net.ResolveUDPAddr(_UDP4, address)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP(_UDP4, nil, addr)
	if err != nil {
		panic(err)
	}

	return &UdpBroadcast{
		conn: conn,
	}
}

func (sock *UdpBroadcast) Close() {
	sock.conn.Close()
}

const (
	maxDatagramSize = 1024
)

type MsgHandler func(int, []byte)

func ListenUdp(address string, handler MsgHandler) {

	addr, err := net.ResolveUDPAddr(_UDP4, address)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenMulticastUDP(_UDP4, nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.SetReadBuffer(maxDatagramSize)



	for {
		buffer := make([]byte, maxDatagramSize)

		numBytes, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}

		handler(numBytes, buffer[:numBytes])
	}
}