package sse

import (
	"context"
	"fmt"
	"io"

	pc "avito/pkg/context"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type Event struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan Message

	// New client connections
	NewClients chan *eventData //chan string

	// Closed client connections
	ClosedClients chan *eventData //chan string

	// Total client connections
	TotalClients map[uuid.UUID]chan Message //chan string

	onCallBackCheckPermissions func(userID uuid.UUID, data interface{}) bool
}

// ClientChan New event messages are broadcast to all registered client connection channels
type ClientChan chan Message

func StartSSE(c *gin.Context) {

	fmt.Println("startSSE")
	v, ok := c.Get("clientChan")
	if !ok {
		logrus.Error("v, ok", v, ok)

		return
	}
	clientChan, ok := v.(ClientChan)
	if !ok {
		logrus.Error("clientChan, ok", clientChan, ok)

		return
	}
	c.Stream(func(w io.Writer) bool {
		// Stream message to client from message channel
		//logrus.Info("Ожидание пакета SSE")

		if msg, ok := <-clientChan; ok {
			//logrus.Info("Пакет SSE отправлен")

			c.SSEvent(string(msg.Event), msg.Data)
			// c.SSEvent("message", msg)
			return true
		}
		return false
	})

}

// NewServerSSE Initialize event and Start procnteessing requests
func NewServerSSE(ctx context.Context) (event *Event) {

	logrus.Info("Сервер SSE запущен")

	event = &Event{
		// Емкость буфера расчина по следующей формуле:
		// (Количество каналов) Х (Количество одновременно запущенных горутин для отправки)
		//
		// Сейчас у нас 2 канала:
		// 1) Обновление сфетофора;
		// 2) Разблокировка плашки в списке ТПН при окончании трехдневного пересчета статистики для нового ТПН
		//
		// Количество одновременно запущенных горутин пересчета статистики равно 10
		// Итого, емкость буфера = 2 x 10 = 20
		Message:       make(chan Message, 20),
		NewClients:    make(chan *eventData),
		ClosedClients: make(chan *eventData),
		TotalClients:  make(map[uuid.UUID]chan Message),
	}

	go event.listen(ctx)

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *Event) listen(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Сервер SSE остановлен...")
			return

		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client.ID] = client.Chan
			logrus.Infof("По SSE подключился клиент: %v. Итого: %d", client.ID, len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client.ID)
			close(client.Chan)
			logrus.Infof("По SSE одключился клиент: %v. Итого: %d", client.ID, len(stream.TotalClients))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			if ch, ok := stream.TotalClients[eventMsg.UserID]; ok {
				ch <- eventMsg
			}

			// for userID, clientMessageChan := range stream.TotalClients {
			// 	// Вызов определенной логики, которая регламентирует отправку по SSE
			// 	res := stream.onCallBackCheckPermissions(userID, eventMsg.Data.(uuid.UUID))

			// 	if res {
			// 		clientMessageChan <- eventMsg
			// 	}

			// }
		}
	}
}

func (stream *Event) ServeHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Info("По SSE установленно соединение с клиентом: ", c.Request.RemoteAddr)

		// Initialize client channel
		clientChan := make(ClientChan)

		data := &eventData{
			ID:   pc.GetUserID(c.Request.Context()),
			Chan: clientChan,
		}

		// Send new connection to event server
		stream.NewClients <- data

		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- data
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

func (stream *Event) OnCallBackCheckPermissions(callback func(userID uuid.UUID, data interface{}) bool) {

	stream.onCallBackCheckPermissions = callback

}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
