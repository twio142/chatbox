package websocket

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if os.Getenv("APP_ENV") == "development" {
			return true
		}
		origin := r.Header.Get("Origin")
		if origin == "" {
			return false
		}
		return origin == "http://" + r.Host || origin == "https://" + r.Host
	},
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}
