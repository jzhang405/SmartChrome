package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

type StreamMessage struct {
	Type      string      `json:"type"`
	Content   string      `json:"content"`
	MessageID string      `json:"message_id"`
	Data      interface{} `json:"data,omitempty"`
}

type StreamManager struct {
	clients map[string]*Client
	mutex   sync.RWMutex
}

type Client struct {
	conn         *websocket.Conn
	conversationID string
	messageID     string
	send         chan StreamMessage
}

func NewStreamManager() *StreamManager {
	return &StreamManager{
		clients: make(map[string]*Client),
	}
}

func (sm *StreamManager) HandleWebSocket(c *gin.Context) {
	conversationID := c.Query("conversationId")
	messageID := c.Query("messageId")

	if conversationID == "" || messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversationId and messageId are required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		conn:          conn,
		conversationID: conversationID,
		messageID:     messageID,
		send:          make(chan StreamMessage, 256),
	}

	clientKey := conversationID + ":" + messageID

	sm.mutex.Lock()
	sm.clients[clientKey] = client
	sm.mutex.Unlock()

	go sm.writePump(client)
	go sm.readPump(client)
}

func (sm *StreamManager) writePump(client *Client) {
	defer func() {
		client.conn.Close()
	}()

	for {
		message, ok := <-client.send
		if !ok {
			client.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := client.conn.WriteJSON(message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
}

func (sm *StreamManager) readPump(client *Client) {
	defer func() {
		clientKey := client.conversationID + ":" + client.messageID
		sm.mutex.Lock()
		delete(sm.clients, clientKey)
		sm.mutex.Unlock()
		client.conn.Close()
	}()

	client.conn.SetReadLimit(512)
	// client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		// client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

func (sm *StreamManager) SendMessage(conversationID, messageID string, message StreamMessage) {
	clientKey := conversationID + ":" + messageID

	sm.mutex.RLock()
	client, exists := sm.clients[clientKey]
	sm.mutex.RUnlock()

	if exists {
		select {
		case client.send <- message:
		default:
			log.Printf("Client send buffer full, dropping message")
		}
	}
}

func (sm *StreamManager) BroadcastToConversation(conversationID string, message StreamMessage) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for clientKey, client := range sm.clients {
		if client.conversationID == conversationID {
			select {
			case client.send <- message:
			default:
				log.Printf("Client send buffer full, dropping message for %s", clientKey)
			}
		}
	}
}

func (sm *StreamManager) SendStreamResponse(ctx context.Context, conversationID, messageID string, content string, isComplete bool) {
	message := StreamMessage{
		Type:      "stream",
		Content:   content,
		MessageID: messageID,
		Data: map[string]interface{}{
			"is_complete": isComplete,
		},
	}

	sm.SendMessage(conversationID, messageID, message)
}

func (sm *StreamManager) SendError(ctx context.Context, conversationID, messageID string, errorMsg string) {
	message := StreamMessage{
		Type:      "error",
		Content:   errorMsg,
		MessageID: messageID,
	}

	sm.SendMessage(conversationID, messageID, message)
}