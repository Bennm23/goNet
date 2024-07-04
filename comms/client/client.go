package main

import (
	"fmt"
	"gosock/protos/netms"
	"time"

	"gosock/comms/multicast"
	"gosock/comms/tcp"

	"google.golang.org/protobuf/proto"
)


func testTcp() {

	go client()

	server := tcp.NewServer(":8897")
	go startTcpServer(server)

	for {
		server.Send([]byte("Server Saying Hello"))
		time.Sleep(time.Second * 2)
	}
}

func startTcpServer(server *tcp.Server) {
	defer server.Close()
	server.Start()

}

func client() {
	client := tcp.NewClient("127.0.0.1:8897", "Client 1")
	defer client.Close()

	go client.Listen(handleTcp)

	for {
		client.Send([]byte("Hello TCP"))
		time.Sleep(time.Second * 1)
	}

}

func handleTcp(bytes []byte) {
	fmt.Println("Got Msg in TCP CLient = ", string(bytes))
}

func main() {

	testTcp()

	// testUdp()
}


func testUdp() {
	fmt.Println("Launching Client")

	test := netms.Test {
		MyInt : 12,
		MyName : "12345",
	}

	fmt.Println(test.String())



	go multicast.ListenUdp("239.0.0.0:8878", handleUdp)

	broadcast := multicast.NewBroadcast("239.0.0.0:8878")
	defer broadcast.Close()

	for {

		broadcast.SendProto(&test)
		time.Sleep(1 * time.Second)
	}
}

func handleUdp(numBytes int, bytes []byte) {

	fmt.Printf("Handling %d Bytes\n", numBytes)
	fmt.Println("Recv = ", bytes)

	msg := &netms.Test{}

	err := proto.Unmarshal(bytes, msg);
	if err != nil {
		panic(err)
	}

	fmt.Println("Proto = ", msg.String())
}