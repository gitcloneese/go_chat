package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

var onlineConns map[string]net.Conn
var messageQue = make(chan string, 100)

func handle(conn net.Conn) {
	buf := make([]byte, 1024)
	defer func(conn net.Conn) {
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		fmt.Println("客户端: ", addr, "退出")
		delete(onlineConns, addr)
		conn.Close()
		for i := range onlineConns {
			fmt.Println("now oline client : ", i)
		}
	}(conn)

	for {
		NumOfReadBytes, err := conn.Read(buf)

		if err != nil && err != io.EOF {
			fmt.Printf("read conn err = %s\n", err)
			break
		}
		if NumOfReadBytes != 0 {
			remoteAddr := conn.RemoteAddr()
			fmt.Print(remoteAddr)
			fmt.Printf(" send message :%s\n", string(buf[:NumOfReadBytes]))
			message := string(buf[:NumOfReadBytes])

			if strings.ToUpper(message) == "EXIT" {
				conn.Write([]byte("exit ack"))
				break
			}
			messageQue <- message
		}

	}

}

func ProcessConsumer() {

	for {
		select {
		case message := <-messageQue:
			content := strings.Split(message, "#")
			if len(content) > 1 {
				addr := content[0]
				lenAddr := len(addr)
				sendMesage := string([]byte(message)[lenAddr+1:])
				addr = strings.Trim(addr, " ")
				if conn, ok := onlineConns[addr]; ok {
					_, err := conn.Write([]byte(sendMesage))
					if err != nil {
						fmt.Println("online conns send failure")
					}
				}
			}
		}
	}

}

func main() {
	onlineConns = make(map[string]net.Conn)
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("listen err := ", err)
		panic(err)
	}

	defer listen.Close()

	fmt.Println("server is waiting!")

	go ProcessConsumer()
	for {

		conn, err := listen.Accept()
		if err != nil {
			fmt.Print("Accept err := ", err)
			continue
		}
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		fmt.Println(addr, " ###  connected server !")
		onlineConns[addr] = conn
		go handle(conn)
	}
}
