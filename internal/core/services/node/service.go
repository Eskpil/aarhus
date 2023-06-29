package node

import (
	"context"
	"github.com/eskpil/aarhus/internal/core"
	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/eskpil/aarhus/pkg/helpers"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func Create(ctx context.Context, s *state.State, input *contracts.NodeCreateInput) (*contracts.Node, error) {
	coll := s.Database.Collection("nodes")

	node := new(contracts.Node)

	node.Id = uuid.New().String()
	node.Observed = false

	node.Location = input.Location
	node.Name = input.Name

	node.Tokens = make([]contracts.NodeToken, 1)
	node.Tokens[0] = contracts.NodeToken{
		Value:  helpers.RandomString(64, helpers.StringTypeAlnum),
		Status: contracts.NodeTokenStatusActive,
	}

	node.Tasks = make([]contracts.Task, 0)

	if _, err := coll.InsertOne(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

func GetAll(ctx context.Context, s *state.State) ([]*contracts.Node, error) {
	coll := s.Database.Collection("nodes")

	var nodes []*contracts.Node

	result, err := coll.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	if err := result.All(ctx, &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

func HandleHeartbeat(ctx context.Context, s *state.State, input *contracts.HeartbeatInput) (*contracts.HeartbeatResponse, error) {
	coll := s.Database.Collection("nodes")

	iden := ctx.Value("identity").(*core.Identity)

	filter := bson.D{
		{"_id", iden.Node.Id},
	}

	update := bson.D{
		{"$set", bson.D{
			{"observed", true},
			{"last_heartbeat", time.Now()},
			{"networking", input.Networking},
		}},
	}

	if _, err := coll.UpdateOne(ctx, filter, update); err != nil {
		return nil, err
	}

	// TODO: Update each server respectively with their server health report which came in with node heartbeat.
	// Missing a server service though.

	response := new(contracts.HeartbeatResponse)

	tasks, err := getPendingTasks(ctx, s, iden.Node.Id)
	if err != nil {
		return nil, err
	}

	response.Tasks = tasks

	return response, nil
}

func CreateTask(ctx context.Context, s *state.State, nodeId string, task *contracts.Task) error {
	coll := s.Database.Collection("nodes")

	filter := bson.D{
		{Key: "_id", Value: nodeId},
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "tasks", Value: task},
		}},
	}

	if _, err := coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func GetById(ctx context.Context, s *state.State, nodeId string) (*contracts.Server, error) {
	coll := s.Database.Collection("nodes")

	server := new(contracts.Server)

	if err := coll.FindOne(ctx, bson.D{{"_id", nodeId}}).Decode(server); err != nil {
		return nil, err
	}

	return server, nil
}

func UpdateTaskStatus(ctx context.Context, s *state.State, nodeId string, taskId string, status contracts.TaskStatus) error {
	coll := s.Database.Collection("nodes")

	filter := bson.D{
		{Key: "_id", Value: nodeId},
		{Key: "tasks.id", Value: taskId},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "tasks.$.status", Value: status},
		}},
	}

	if _, err := coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func getPendingTasks(ctx context.Context, s *state.State, nodeId string) ([]*contracts.Task, error) {
	coll := s.Database.Collection("nodes")

	node := new(contracts.Node)
	if err := coll.FindOne(ctx, bson.D{{Key: "_id", Value: nodeId}}).Decode(node); err != nil {
		return nil, err
	}

	var tasks []*contracts.Task

	for _, task := range node.Tasks {
		if task.Status == contracts.TaskStatusPending || task.Status == contracts.TaskStatusAborted {
			tasks = append(tasks, &task)

			if err := UpdateTaskStatus(ctx, s, nodeId, task.Id, contracts.TaskStatusScheduled); err != nil {
				return nil, err
			}
		}
	}

	return tasks, nil
}
