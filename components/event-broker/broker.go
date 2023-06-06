package EventBroker

import (
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/soulteary/sparrow/internal/define"
)

type Event struct {
	ConversationID  string
	ParentMessageID string
	Payload         any
}

type EventChan chan Event

type Broker struct {
	ID         int
	Event      EventChan
	Connect    chan EventChan
	Disconnect chan EventChan
	Clients    map[EventChan]struct{}
}

func NewBroker(id int) (broker *Broker) {
	return &Broker{
		ID:         id,
		Event:      make(EventChan, 1),
		Connect:    make(chan EventChan),
		Disconnect: make(chan EventChan),
		Clients:    make(map[EventChan]struct{}),
	}
}

func (broker *Broker) Listen() {
	for {
		select {
		case s := <-broker.Connect:
			broker.Clients[s] = struct{}{}
			log.Println("New Connection, number of client connections:", len(broker.Clients))
		case s := <-broker.Disconnect:
			delete(broker.Clients, s)
			log.Println("Disconnect, number of client connections:", len(broker.Clients))
		case event := <-broker.Event:
			for client := range broker.Clients {
				select {
				case client <- event:
				case <-time.After(define.DEFAULT_EVENT_BROKER_PATIENCE):
					log.Panicln("Ignore client connection.")
				}
			}
		}
	}
}

func (broker *Broker) Serve(c *gin.Context, messageChan EventChan) {
	conversationID := c.Request.Header.Get("x-conversation-id")
	parentMessageID := c.Request.Header.Get("x-parent-message-id")

	log.Println("Requested topic:", conversationID, parentMessageID)

	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Transfer-Encoding", "chunked")
	c.Status(200)

	broker.Connect <- messageChan
	defer func() {
		broker.Disconnect <- messageChan
	}()

	c.Stream(func(w io.Writer) bool {
		event := <-messageChan
		if event.ConversationID != conversationID || event.ParentMessageID != parentMessageID {
			return false
		}
		if IsLastMessage(event.Payload) {
			c.SSEvent("", event.Payload)
			broker.Disconnect <- messageChan
			c.Abort()
			c.Writer.CloseNotify()
			FreePool(broker.ID)
		} else {
			c.SSEvent("", event.Payload)
		}
		c.Writer.Flush()
		return true
	})
}
