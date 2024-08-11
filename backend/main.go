package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/rs/cors"
	"github.com/twio142/realtime-chat-go-react/pkg/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	if os.Getenv("APP_ENV") == "production" {
		fs := http.FileServer(http.Dir("./build"))
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			w.Header().Set("Content-Security-Policy", "script-src 'self'")
			w.Header().Set("Content-Security-Policy", "connect-src 'self'")
			fs.ServeHTTP(w, r)
		}))
	}

	corsWrapper := cors.New(cors.Options{}).Handler
	if os.Getenv("APP_ENV") == "development" {
		corsWrapper = cors.AllowAll().Handler
	}

	http.Handle("/upload", corsWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(128 << 20)
		if err != nil {
			http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error getting the file from the form: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filename := header.Filename
		for {
			if _, err := os.Stat("./uploads/" + filename); os.IsNotExist(err) {
				break
			}

			root := filepath.Base(filename)
			ext := filepath.Ext(filename)
			root = root[:len(root)-len(ext)]

			re := regexp.MustCompile(`( \(\d+\))?$`)
			root = re.ReplaceAllStringFunc(root, func(s string) string {
				if s == "" {
					return " (1)"
				}
				num, _ := strconv.Atoi(s[2 : len(s)-1])
				return fmt.Sprintf(" (%d)", num+1)
			})

			filename = root + ext
		}

		dst, err := os.Create("./uploads/" + filename)
		if err != nil {
			http.Error(w, "Error creating the file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error copying the file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fmt.Sprintf(`{"fileName": "%s", "fileURL": "/u/%s"}`, header.Filename, filename)))
	})))

	http.Handle("/u/", corsWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := "./uploads/" + r.URL.Path[3:]

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, filePath)
	})))
}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	setupRoutes()
	address := os.Getenv("PORT")
	if address == "" {
		address = "0.0.0.0:8080"
	}
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
