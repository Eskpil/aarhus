package ticket

import (
	"context"
	"github.com/eskpil/aarhus/internal/core"
	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func FindTicket(ctx context.Context, s *state.State, ticketId string) (*contracts.Ticket, error) {
	coll := s.Database.Collection("tickets")

	ticket := new(contracts.Ticket)

	if err := coll.FindOne(ctx, bson.D{{"_id", ticketId}}).Decode(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func Create(ctx context.Context, s *state.State, input *contracts.TicketInput) (*contracts.Ticket, error) {
	coll := s.Database.Collection("tickets")

	identity := ctx.Value("identity").(*core.Identity)

	ticket := new(contracts.Ticket)

	ticket.Id = uuid.New().String()
	ticket.ServerId = input.ServerId
	ticket.UserId = identity.User.Id
	ticket.CreatedAt = time.Now()
	ticket.UpdatedAt = time.Now()

	if _, err := coll.InsertOne(ctx, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}
