package main

import (
	"log"
	"net"

	"github.com/ddimm/crash-course/chat"
)

func main() {

	userChan := make(chan *chat.ChatUser)
	receiveChan := make(chan *chat.ChatMessage)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	chat.ListenerHandler(&listener, userChan, receiveChan)
}
