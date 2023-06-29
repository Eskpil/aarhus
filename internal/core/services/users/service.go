package users

import (
	"context"
	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/pkg/contracts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Replace(ctx context.Context, s *state.State, user *contracts.User) error {
	coll := s.Database.Collection("users")

	filter := bson.D{{"_id", user.Id}}

	opts := new(options.ReplaceOptions)
	opts.SetUpsert(true)

	if _, err := coll.ReplaceOne(ctx, filter, user, opts); err != nil {
		return err
	}

	return nil
}

func FindById(ctx context.Context, s *state.State, uid string) (*contracts.User, error) {
	coll := s.Database.Collection("users")

	filter := bson.D{{"_id", uid}}

	user := new(contracts.User)
	if err := coll.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func FindByEmail(ctx context.Context, s *state.State, email string) (*contracts.User, error) {
	coll := s.Database.Collection("users")

	filter := bson.D{{"email", email}}

	user := new(contracts.User)
	if err := coll.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}
