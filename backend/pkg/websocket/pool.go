package websocket

import (
	"fmt"
	"os"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) getClientNames() map[string]bool {
	names := make(map[string]bool)
	for client := range pool.Clients {
		names[client.ID] = true
	}
	return names
}

func clearUploads() {
	files, err := os.ReadDir("./uploads")
	if err != nil {
		fmt.Println("Error reading uploads directory: ", err)
		return
	}
	for _, file := range files {
		err := os.RemoveAll("./uploads/" + file.Name())
		if err != nil {
			fmt.Println("Error removing file: ", err)
		}
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case clientJoined := <-pool.Register:
			currentNames := pool.getClientNames()
			for {
				newName := RandomName()
				if !currentNames[newName] {
					clientJoined.ID = newName
					break
				}
			}
			pool.Clients[clientJoined] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				if clientJoined.ID == client.ID {
					text := fmt.Sprintf("Welcome, %s.", client.ID)
					client.Conn.WriteJSON(Message{Text: text})
					text = fmt.Sprintf("There %s online.", func() string {
						if len(pool.Clients) == 1 {
							return "is 1 person"
						} else {
							return fmt.Sprintf("are %d people", len(pool.Clients))
						}
					}())
					client.Conn.WriteJSON(Message{Text: text})
				} else {
					text := fmt.Sprintf("%s entered the room.", clientJoined.ID)
					client.Conn.WriteJSON(Message{Text: text})
				}
			}
		case clientLeft := <-pool.Unregister:
			delete(pool.Clients, clientLeft)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				text := fmt.Sprintf("%s left the room.", clientLeft.ID)
				client.Conn.WriteJSON(Message{Text: text})
			}
			if len(pool.Clients) == 0 {
				fmt.Println("No clients left, removing uploads")
				clearUploads()
			}
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client := range pool.Clients {
				text := fmt.Sprintf("%s: %s", message.Client.ID, message.Text)
				_type := 2
				if message.Client.ID == client.ID {
					_type = 1
				}
				client.Conn.WriteJSON(Message{Text: text, FileName: message.FileName, FileURL: message.FileURL, Type: _type})
			}
		}
	}
}
