package main

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	/**
	websocketUpgrader is used to upgrade incomming HTTP requests into a persitent websocket connection
	*/
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

type Manager struct {
	clients ClientList

	// Using a syncMutex here to be able to lcok state before editing clients
	// Could also use Channels to block
	sync.RWMutex
	// handlers are functions that are used to handle Events
	handlers map[string]EventHandler
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request, gameCode string) {
	// Begin by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
		return
	}

	// Create New Client
	client := NewClient(conn, userID, gameCode)
	// Add the newly created client to the manager
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clients[client] = true
}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Check if Client exists, then delete it
	if _, ok := m.clients[client]; ok {
		// close connection
		client.connection.Close()
		// remove
		delete(m.clients, client)
	}
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendAnswer] = SendAnswerHandler
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}
