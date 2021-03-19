package chat

import (
	"net"
)

// Types and function with capital letter names are exported
// everything else is private to the package

// ChatUser represents a connected user to the chat system
type ChatUser struct {
	Name   string
	Conn   *net.Conn
	Sender chan ChatMessage
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Person  *ChatUser
	Message string
}

func messageSender(userChan chan *ChatUser, receiveChan chan *ChatMessage) {

}

// ListenerHandler sets up the messageSender thread and handles
// incoming messages
func ListenerHandler(listener *net.Listener, userChan chan *ChatUser, receiveChan chan *ChatMessage) {

}
