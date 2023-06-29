package server

import (
	"context"
	"github.com/eskpil/aarhus/internal/core"
	nodeService "github.com/eskpil/aarhus/internal/core/services/node"
	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func Create(ctx context.Context, s *state.State, input *contracts.ServerCreateInput) (*contracts.Server, error) {
	coll := s.Database.Collection("servers")

	server := new(contracts.Server)

	server.Id = uuid.New().String()
	server.Name = input.Name
	server.NodeRef = input.NodeRef
	server.PortForwards = input.PortForwards

	// TODO: Verify the server allocation against the resources of the node.
	server.Allocation = input.Allocation
	server.Type = input.Type
	server.MinecraftServer.ServerJar = input.MinecraftServerCreateInput.ServerJar

	if _, err := coll.InsertOne(ctx, server); err != nil {
		return nil, err
	}

	task := new(contracts.Task)

	task.Id = uuid.New().String()
	task.Type = contracts.TaskTypeCreateServer
	task.Status = contracts.TaskStatusPending
	task.TaskCreateServer = contracts.TaskCreateServer{
		ServerId: server.Id,
	}

	if err := nodeService.CreateTask(ctx, s, server.NodeRef, task); err != nil {
		return nil, err
	}

	return server, nil
}

//func GetAll(ctx context.Context, s *state.State) ([]*contracts.Server, error) {
//	coll := s.Database.Collection("servers")
//
//	result, err := coll.Find(ctx, bson.D{{}})
//	if err != nil {
//		return nil, err
//	}
//
//	var servers []*contracts.Server
//	if err := result.Decode(&servers); err != nil {
//		return nil, err
//	}
//
//	return servers, nil
//}

func GetAll(ctx context.Context, s *state.State) ([]*contracts.Server, error) {
	coll := s.Database.Collection("servers")

	var servers []*contracts.Server

	filter := bson.D{{}}

	identity := ctx.Value("identity").(*core.Identity)

	if identity.Type == core.IdentityTypeNode {
		filter = append(filter, bson.E{Key: "node_ref", Value: identity.Node.Id})
	}

	result, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := result.All(ctx, &servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func GetById(ctx context.Context, s *state.State, serverId string) (*contracts.Server, error) {
	coll := s.Database.Collection("servers")

	server := new(contracts.Server)
	if err := coll.FindOne(ctx, bson.D{{"_id", serverId}}).Decode(server); err != nil {
		return nil, err
	}

	return server, nil
}
