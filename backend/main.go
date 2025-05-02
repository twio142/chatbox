package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

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

	http.Handle("/upload/chunk", corsWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form with reasonable chunk size
		if err := r.ParseMultipartForm(8 << 20); err != nil { // 8MB per chunk
			http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get chunk information
		chunkNumber, err := strconv.Atoi(r.FormValue("chunkNumber"))
		if err != nil {
			http.Error(w, "Invalid chunk number", http.StatusBadRequest)
			return
		}

		totalChunks, err := strconv.Atoi(r.FormValue("totalChunks"))
		if err != nil {
			http.Error(w, "Invalid total chunks", http.StatusBadRequest)
			return
		}

		filename := r.FormValue("filename")
		if filename == "" {
			http.Error(w, "Filename required", http.StatusBadRequest)
			return
		}

		// Create unique ID for this file or use provided one
		fileID := r.FormValue("fileID")
		if fileID == "" {
			http.Error(w, "File ID required", http.StatusBadRequest)
			return
		}

		// Get the actual chunk data
		file, _, err := r.FormFile("chunk")
		if err != nil {
			http.Error(w, "Error getting chunk from form: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Create uploads directory if it doesn't exist
		chunksDir := "./uploads/chunks/" + fileID
		if _, e := os.Stat(chunksDir); os.IsNotExist(e) {
			e = os.MkdirAll(chunksDir, 0755)
			if e != nil {
				http.Error(w, "Error creating chunks directory: "+e.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Write this chunk to a temporary file
		chunkPath := filepath.Join(chunksDir, fmt.Sprintf("chunk-%d", chunkNumber))
		dst, err := os.Create(chunkPath)
		if err != nil {
			http.Error(w, "Error creating chunk file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving chunk: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if all chunks are received
		if chunkNumber == totalChunks-1 {
			// All chunks received, merge them
			finalPath := filepath.Join("./uploads", filename)

			// Make sure the uploads directory exists
			if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
				err = os.Mkdir("./uploads", 0755)
				if err != nil {
					http.Error(w, "Error creating uploads directory: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			// Create the final file
			finalFile, err := os.Create(finalPath)
			if err != nil {
				http.Error(w, "Error creating final file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer finalFile.Close()

			// Append each chunk to the final file
			for i := range totalChunks {
				chunkPath := filepath.Join(chunksDir, fmt.Sprintf("chunk-%d", i))
				chunkData, err := os.Open(chunkPath)
				if err != nil {
					http.Error(w, "Error opening chunk: "+err.Error(), http.StatusInternalServerError)
					return
				}

				_, err = io.Copy(finalFile, chunkData)
				chunkData.Close()
				if err != nil {
					http.Error(w, "Error appending chunk: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			// Remove chunks directory
			os.RemoveAll(chunksDir)

			// Return success response with file URL
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"fileName": "%s", "fileURL": "/u/%s", "status": "complete"}`, filename, filename)
		} else {
			// Return progress response
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status": "chunk-received", "chunkNumber": %d, "totalChunks": %d}`, chunkNumber, totalChunks)
		}
	})))

	http.Handle("/upload", corsWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 30); err != nil {
			http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error getting the file from the form: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		if _, e := os.Stat("./uploads"); os.IsNotExist(e) {
			os.Mkdir("./uploads", 0755)
		}

		filename := header.Filename
		for {
			if _, e := os.Stat("./uploads/" + filename); os.IsNotExist(e) {
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
		_, _ = fmt.Fprintf(w, `{"fileName": "%s", "fileURL": "/u/%s"}`, header.Filename, filename)
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port " + port)
	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Minute,
		WriteTimeout:      30 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
