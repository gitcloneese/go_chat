package main

import (
	"fmt"
	"time"
)

var obchan chan interface{}

type ob struct {
	name string
}

func consumer() {
	for {
		select {
		case message := <-obchan:
			fmt.Println(message)

		}
	}

}
func main() {
	go consumer()
	obchan = make(chan interface{}, 10)
	test := ob{name: "xiaoxiao"}
	obchan <- test
	test.name = "dada"
	fmt.Println(test)
	time.Sleep(time.Second * 10)

}
