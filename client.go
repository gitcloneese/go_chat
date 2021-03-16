package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("conn err := ", err)
		panic(err)
	}

	defer conn.Close()

	go MessageSend(conn)

	buf := make([]byte, 1024)

	for {
		length, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read err := ", err)
		}

		if length != 0 {
			message := buf[:length]
			if strings.ToUpper(string(message)) == "EXIT ACK" {
				fmt.Println("客户端退出， 欢迎使用北风网聊天系统")
				os.Exit(0)
			}
			fmt.Println("read message from server :", string(buf[:length]))
		}
	}

}

func MessageSend(conn net.Conn) {
	defer conn.Close()
	var input string

	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)

		if strings.ToUpper(input) == "EXIT" {
			_, err1 := conn.Write(data)
			if err1 != nil {
				fmt.Println("exit err")
			}
		} else {
			_, err := conn.Write([]byte(input))
			if err != nil {
				fmt.Println("client connect err : ", err)
			}
		}

	}

}
