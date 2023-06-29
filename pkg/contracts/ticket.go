package contracts

import "time"

type TicketStatus int64

const (
	TicketStatusUnknown TicketStatus = iota
	TicketStatusOpen
	TicketStatusClosed
)

type Ticket struct {
	Id string `json:"id" bson:"_id"`

	UserId   string `json:"user_id" bson:"user_id"`
	ServerId string `json:"server_id" bson:"server_id"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type TicketInput struct {
	ServerId string `json:"server_id"`
}
