package state

import (
	"context"
	"fmt"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

type State struct {
	Nodes []*Node

	Oauth2Config *oauth2.Config

	Database *mongo.Database
}

func New(ctx context.Context) (*State, error) {
	state := new(State)

	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s",
		viper.GetString("MONGODB_USER"),
		viper.GetString("MONGODB_PASSWORD"),
		viper.GetString("MONGODB_HOST"),
		viper.GetString("MONGODB_PORT"),
		viper.GetString("MONGODB_DATABASE"),
	)

	slog.Infoc(ctx, "connecting with mongodb", slog.String("uri", uri))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	state.Database = client.Database(viper.GetString("MONGODB_DATABASE"))

	state.Oauth2Config = &oauth2.Config{
		RedirectURL: "http://localhost:8080/v1/auth/callback/",
		// This next 2 lines must be edited before running this.
		ClientID:     viper.GetString("DISCORD_CLIENT_ID"),
		ClientSecret: viper.GetString("DISCORD_CLIENT_SECRET"),
		Scopes:       []string{"identity", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://discord.com/api/oauth2/authorize",
			TokenURL:  "https://discord.com/api/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	slog.Info("initialized state")

	return state, nil
}

func (s *State) FindNode(nodeId string) (*Node, error) {
	for _, node := range s.Nodes {
		if node.Id == nodeId {
			return node, nil
		}
	}

	return nil, fmt.Errorf("could not find node")
}
