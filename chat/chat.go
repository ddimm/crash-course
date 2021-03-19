package chat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
	// we'll use a map to keep track of users on the service
	// we'll update this map within this function, so it should be safe
	users := make(map[string]*ChatUser)
	for {
		// the select state is like a switch statement for channels
		select {
		case newMessage := <-receiveChan:
			// do something

			for name, user := range users {
				if name != newMessage.Person.Name {
					// might have timeout problems
					user.Sender <- *newMessage

				}
			}

		case newUser := <-userChan:
			// do something with new users
			_, ok := users[newUser.Name]
			if ok {
				delete(users, newUser.Name)
				for name, user := range users {
					if name != newUser.Name {
						user.Sender <- ChatMessage{Person: newUser, Message: fmt.Sprintf("%s has left\n", newUser.Name)}

					}
				}
			} else {
				users[newUser.Name] = newUser
				for name, user := range users {
					if name != newUser.Name {
						user.Sender <- ChatMessage{Person: newUser, Message: fmt.Sprintf("%s has joined\n", newUser.Name)}

					}
				}

			}

		}
	}
}

// ListenerHandler sets up the messageSender thread and handles
// incoming messages
func ListenerHandler(listener *net.Listener, userChan chan *ChatUser, receiveChan chan *ChatMessage) {
	// all we need to do to get concurrent
	go messageSender(userChan, receiveChan)
	for {
		conn, err := (*listener).Accept()
		// common way to check for errors
		if err != nil {
			log.Println(err)
			continue
		}
		// we'll use an anon. function
		go func() {
			// use defer to close the connection later
			defer conn.Close()
			// get a zeroed user with new
			newUser := new(ChatUser)
			newUser.Conn = &conn
			// we'll use bufio to read from the connection
			buf := bufio.NewReader(conn)
			name, err := buf.ReadString('\n')
			if err != nil {
				log.Println(err)
				return
			}
			newUser.Name = strings.TrimSpace(name)
			newUser.Sender = make(chan ChatMessage)
			// we'll also defer close the sender channel
			defer close(newUser.Sender)
			userChan <- newUser
			go func() {
				for m := range newUser.Sender {
					io.WriteString(conn, fmt.Sprintf("%s: %s", m.Person.Name, m.Message))
				}
			}()
			for {
				// we'll run in this loop until we read EOF from the connection
				var m ChatMessage
				m.Person = newUser
				newMessage, err := buf.ReadString('\n')
				if err != nil {
					log.Println(err)
					break
				}
				m.Message = newMessage
				receiveChan <- &m
			}
			userChan <- newUser
		}()

	}
}
