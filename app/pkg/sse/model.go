package sse

import "github.com/gofrs/uuid"

type eventData struct {
	ID   uuid.UUID
	Chan chan Message
}

type Message struct {
	UserID uuid.UUID
	Event  EventName
	Data   interface{}
}

type EventName string

const (
	EventTCMs EventName = "tcms"
)
