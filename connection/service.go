package connection

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type ConnectionService struct {
	Connections      map[string]*websocket.Conn
	ConnectionsMutex sync.Mutex
}

var connectionServiceInstance *ConnectionService
var once sync.Once

func ProvideConnectionService() *ConnectionService {
	log.Println("ProvideConnectionService")
	return GetConnectionServiceInstance()
}

func GetConnectionServiceInstance() *ConnectionService {
	once.Do(func() {
		connectionServiceInstance = &ConnectionService{
			Connections:      make(map[string]*websocket.Conn),
			ConnectionsMutex: sync.Mutex{},
		}
		log.Println("🔌 Successfully connected to Connection Service")
	})
	return connectionServiceInstance
}

func (e *ConnectionService) PrintAllConnections() {
	for key, _ := range e.Connections {
		fmt.Println(key)
	}
}

func (e *ConnectionService) AddConnection(key string, conn *websocket.Conn) {
	e.ConnectionsMutex.Lock()
	defer e.ConnectionsMutex.Unlock()
	e.Connections[key] = conn
}

func (e *ConnectionService) GetConnection(key string) (*websocket.Conn, bool) {
	e.ConnectionsMutex.Lock()
	defer e.ConnectionsMutex.Unlock()
	conn, ok := e.Connections[key]
	return conn, ok
}

func (e *ConnectionService) RemoveConnection(key string) {
	e.ConnectionsMutex.Lock()
	defer e.ConnectionsMutex.Unlock()
	delete(e.Connections, key)
	log.Println("❌ Successfully removed connection")
}

func (e *ConnectionService) UpdateConnection(key string, conn *websocket.Conn) {
	e.ConnectionsMutex.Lock()
	defer e.ConnectionsMutex.Unlock()
	e.Connections[key] = conn
}

func (e *ConnectionService) GetConnections() map[string]*websocket.Conn {
	e.ConnectionsMutex.Lock()
	defer e.ConnectionsMutex.Unlock()
	return e.Connections
}

// send message to all connections
func (e *ConnectionService) SendToAll(message string) {
	e.ConnectionsMutex.Lock()
	defer e.ConnectionsMutex.Unlock()
	for _, conn := range e.Connections {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("write:", err)
		}
	}
}
