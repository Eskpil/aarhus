package node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eskpil/aarhus/internal/node/scheduler"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"os"
)

type State struct {
	NodeToken string

	Scheduler *scheduler.Scheduler

	Client *resty.Client

	Upgrader *websocket.Upgrader

	Sessions map[string][]*websocket.Conn
}

func New() (*State, error) {
	state := new(State)

	if err := state.LoadNodeToken(); err != nil {
		return nil, err
	}

	sched, err := scheduler.New()
	if err != nil {
		return nil, err
	}
	state.Scheduler = sched

	slog.Info("state", slog.String("node_token", state.NodeToken))

	state.Client = resty.New()
	state.Client.SetHeader("X-Aarhus-Node-Token", state.NodeToken)

	state.Upgrader = new(websocket.Upgrader)

	return state, nil
}

func (s *State) LoadNodeToken() error {
	path := fmt.Sprintf("%s/node_token", viper.GetString("AARHUS_NODE_DATA"))

	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	s.NodeToken = string(bytes)[:64]

	return nil
}

func (s *State) Heartbeat(ctx context.Context) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	body := &contracts.HeartbeatInput{
		Networking: contracts.CommonNetworking{
			Hostname: hostname,
		},
		HealthReports: []contracts.ServerHealthReport{},
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := s.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(bytes).
		Post(fmt.Sprintf("%s/v1/heartbeat/", viper.GetString("CORE")))
	if err != nil {
		return err
	}

	response := new(contracts.HeartbeatResponse)
	if err := json.Unmarshal(resp.Body(), response); err != nil {
		return err
	}

	for _, task := range response.Tasks {
		if err := s.processTask(ctx, task); err != nil {
			slog.Errorc(ctx, "could not process task", err)
		}
	}

	return nil
}

func (s *State) processTask(ctx context.Context, task *contracts.Task) error {

	switch task.Type {
	case contracts.TaskTypeCreateServer:
		{
			resp, err := s.Client.R().
				SetContext(ctx).
				Get(fmt.Sprintf(`%s/v1/servers/%s`, viper.GetString("CORE"), task.TaskCreateServer.ServerId))
			if err != nil {
				return err
			}

			server := new(contracts.Server)
			if err := json.Unmarshal(resp.Body(), server); err != nil {
				return err
			}

			input := new(scheduler.ServerInput)

			input.Id = server.Id
			input.Name = server.Name
			input.ServerJar = server.MinecraftServer.ServerJar
			input.Ram = server.Allocation.Ram

			_, err = s.Scheduler.ScheduleServer(ctx, input)
			if err != nil {
				_, err := s.Client.R().
					SetContext(ctx).
					SetBody(map[string]int64{"status": int64(contracts.TaskStatusAborted)}).
					Put(fmt.Sprintf("%s/v1/@me/tasks/%s/status", viper.GetString("CORE"), task.Id))
				if err != nil {
					return err
				}

				return err
			}
		}
	default:
		return fmt.Errorf("unknown task type: %d", task.Type)
	}

	return nil
}

func (s *State) ValidateTicket(ctx context.Context, ticketId string) (*contracts.Ticket, error) {
	resp, err := s.Client.R().
		SetContext(ctx).
		Get(fmt.Sprintf("%s/v1/ticket/%s/", viper.GetString("CORE"), ticketId))
	if err != nil {
		return nil, err
	}

	ticket := new(contracts.Ticket)
	if err := json.Unmarshal(resp.Body(), &ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *State) Start() error {
	return nil
}
