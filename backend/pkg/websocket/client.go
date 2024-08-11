package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"encoding/json"
	// "sync"

	"github.com/gorilla/websocket"
)

var names = []string{"Hippo", "Capybara", "Giraffe", "Meerkat", "Otter", "Penguin", "Panda", "Lion", "Tiger", "Elephant", "Zebra", "Kangaroo", "Koala", "Monkey", "Sloth", "Snake", "Turtle", "Frog", "Lizard", "Chameleon", "Alligator", "Crocodile", "Iguana", "Dragon", "Dinosaur", "Unicorn", "Pegasus", "Phoenix", "Griffin", "Mermaid", "Siren", "Centaur", "Minotaur", "Satyr", "Cyclops", "Goblin", "Orc", "Troll", "Gnome", "Dwarf", "Elf", "Fairy", "Angel", "Demon", "Vampire", "Werewolf", "Ghost", "Zombie", "Mummy", "Witch", "Wizard", "Warlock", "Sorcerer", "Necromancer", "Shaman", "Druid", "Priest", "Monk", "Paladin", "Barbarian", "Ranger", "Rogue", "Bard", "Fighter", "Cleric", "Rogue", "Sorcerer", "Warlock", "Wizard", "Bard", "Druid", "Monk", "Paladin", "Ranger", "Cleric", "Fighter", "Barbarian", "Rogue", "Sorcerer", "Warlock", "Wizard", "Bard", "Druid", "Monk", "Paladin", "Ranger", "Cleric", "Fighter", "Barbarian", "Rogue", "Sorcerer", "Warlock", "Wizard", "Bard", "Druid", "Monk", "Paladin", "Ranger", "Cleric", "Fighter", "Barbarian", "Rogue", "Sorcerer", "Warlock", "Wizard", "Bard", "Druid", "Monk", "Paladin", "Ranger", "Cleric", "Fighter", "Barbarian"}

func RandomName() string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return names[r.Intn(len(names))]
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type      int      `json:"type"`
	Text      string   `json:"text"`
	Client    *Client  `json:"client"`
	FileName   string  `json:"fileName,omitempty"`
	FileURL    string  `json:"fileURL,omitempty"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{}
		err = json.Unmarshal([]byte(string(p)), &message)
		if err != nil {
			log.Printf("Error parsing JSON: %v", err)
			return
		}
		message.Type = messageType
		message.Client = c
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
