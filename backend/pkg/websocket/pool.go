package websocket

import (
	"fmt"
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

func (pool *Pool) getClientNames() []string {
	names := make([]string, 0, len(pool.Clients))
	for client := range pool.Clients {
		names = append(names, client.ID)
	}
	return names
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func (pool *Pool) Start() {
	for {
		select {
		case clientJoined := <-pool.Register:
			for {
				newName := RandomName()
				names := pool.getClientNames()
				if !contains(names, newName) {
					clientJoined.ID = newName
					break
				}
			}
			pool.Clients[clientJoined] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				if clientJoined.ID == client.ID {
					text := fmt.Sprintf("Welcome, %s.", client.ID)
					client.Conn.WriteJSON(Message{Text: text})
					text = fmt.Sprintf("There %s online.", func() string {
						if len(pool.Clients) == 1 {
							return "is 1 person"
						} else {
							return fmt.Sprintf("are %d people.", len(pool.Clients))
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
			for client, _ := range pool.Clients {
				text := fmt.Sprintf("%s left the room.", clientLeft.ID)
				client.Conn.WriteJSON(Message{Text: text})
			}
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				text := fmt.Sprintf("%s: %s", message.Client.ID, message.Text)
				client.Conn.WriteJSON(Message{Text: text, FileName: message.FileName, FileURL: message.FileURL})
			}
		}
	}
}
