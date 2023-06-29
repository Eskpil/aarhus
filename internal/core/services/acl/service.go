package acl

import (
	"context"
	"errors"
	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/pkg/contracts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyNodeIdentity(ctx context.Context, s *state.State, nodeToken string) (*contracts.Node, error) {
	coll := s.Database.Collection("nodes")

	node := new(contracts.Node)
	if err := coll.FindOne(
		ctx,
		bson.D{
			{"tokens.value", nodeToken},
			{"tokens.status", contracts.NodeTokenStatusActive},
		}).Decode(node); err != nil {
		return nil, err
	}

	return node, nil
}

func CheckUserAccess(ctx context.Context, s *state.State, user *contracts.User) (bool, error) {
	coll := s.Database.Collection("access")

	entry := new(contracts.AccessEntry)

	if err := coll.FindOne(ctx, bson.D{{"email", user.Email}}).Decode(&entry); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
