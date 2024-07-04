package tcp

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Server struct {
	connections []net.Conn;
	port string;
}

func NewServer(port string) *Server{

	return &Server{
		connections: make([]net.Conn, 0),
		port: port,
	}


}

func (server *Server) Send(bytes []byte) {


	length := make([]byte, 4);
	binary.LittleEndian.PutUint32(length, uint32(len(bytes)))

	for _, conn := range server.connections {

		fmt.Println("Server Sending ", string(bytes))
		conn.Write(length)
		conn.Write(bytes)
	}

}
func (server *Server) Close() {

	for _, conn := range server.connections {
		conn.Close()
	}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp4", server.port)

	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		server.connections = append(server.connections, conn)

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	for {
		lengthBytes := make([]byte, 4)

		_, err := conn.Read(lengthBytes)
		if err != nil {
			panic(err)
		}

		length := binary.LittleEndian.Uint32(lengthBytes)

		content := make([]byte, length)

		_, err = conn.Read(content)

		if err != nil {
			panic(err)
		}

		fmt.Println("Server Conn Read Msg = ", string(content))
	}

}

type Client struct {
	connection net.Conn;
}
func NewClient(addr, name string) *Client{
	connection, err := net.Dial("tcp", addr)

	if err != nil {
		panic(err)
	}

	return &Client{connection: connection}
}

func (client Client) Send(bytes []byte) {

	length := make([]byte, 4);
	binary.LittleEndian.PutUint32(length, uint32(len(bytes)))

	fmt.Println("Client Sending ", string(bytes))
	client.connection.Write(length)

	client.connection.Write(bytes)
}

func (client Client) Close() {
	client.connection.Close()
}

type Callback func([]byte)

func (client Client) Listen(handler Callback) {

	for {

		lengthBytes := make([]byte, 4)

		_, err := client.connection.Read(lengthBytes)
		if err != nil {
			panic(err)
		}

		length := binary.LittleEndian.Uint32(lengthBytes)

		content := make([]byte, length)

		_, err = client.connection.Read(content)

		if err != nil {
			panic(err)
		}
		handler(content)

	}

}