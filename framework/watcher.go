package framework

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type LiveReload struct {
	upgrader  websocket.Upgrader
	clients   map[*websocket.Conn]bool
	clientsMu sync.Mutex
	watcher   *fsnotify.Watcher
}

func NewLiveReload() (*LiveReload, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	lr := &LiveReload{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		clients: make(map[*websocket.Conn]bool),
		watcher: watcher,
	}

	go lr.watchLoop()
	log.Printf("LiveReload initialized")
	return lr, nil
}

func (lr *LiveReload) watchLoop() {
	for {
		select {
		case event, ok := <-lr.watcher.Events:
			if !ok {
				return
			}
			log.Printf("File event detected: %v", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("File modified: %s", event.Name)
				lr.broadcastReload()
			}
		case err, ok := <-lr.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error watching files: %v", err)
		}
	}
}

func (lr *LiveReload) broadcastReload() {
	lr.clientsMu.Lock()
	defer lr.clientsMu.Unlock()

	log.Printf("Broadcasting reload to %d clients", len(lr.clients))
	for client := range lr.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
		if err != nil {
			log.Printf("Error sending reload message: %v", err)
			client.Close()
			delete(lr.clients, client)
		} else {
			log.Printf("Reload message sent successfully")
		}
	}
}

func (lr *LiveReload) HandleWebSocket(c echo.Context) error {
	log.Printf("New WebSocket connection request from %s", c.Request().RemoteAddr)
	ws, err := lr.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return err
	}

	// Configure WebSocket connection
	ws.SetReadLimit(512) // Limit size of incoming messages
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	lr.clientsMu.Lock()
	lr.clients[ws] = true
	clientCount := len(lr.clients)
	lr.clientsMu.Unlock()
	log.Printf("New WebSocket client connected. Total clients: %d", clientCount)

	// Start ping-pong routine
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					log.Printf("Failed to send ping: %v", err)
					return
				}
			}
		}
	}()

	// Keep connection alive and handle incoming messages
	for {
		messageType, _, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		if messageType == websocket.PingMessage {
			if err := ws.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				log.Printf("Failed to send pong: %v", err)
				break
			}
		}
	}

	// Cleanup when connection closes
	lr.clientsMu.Lock()
	delete(lr.clients, ws)
	clientCount = len(lr.clients)
	lr.clientsMu.Unlock()
	ws.Close()
	log.Printf("WebSocket client disconnected. Remaining clients: %d", clientCount)
	return nil
}

func (lr *LiveReload) WatchDir(dir string) error {
	log.Printf("Starting to watch directory: %s", dir)
	return lr.watcher.Add(dir)
}

func (lr *LiveReload) Close() error {
	lr.clientsMu.Lock()
	for client := range lr.clients {
		client.Close()
	}
	lr.clientsMu.Unlock()
	return lr.watcher.Close()
}
