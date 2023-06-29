package scheduler

import (
	"context"
	"fmt"
	"github.com/spf13/viper"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/eskpil/aarhus/pkg/utils"
)

type ServerInput struct {
	Id string

	Name      string
	ServerJar string
	Ram       uint64
	VCpu      float64
}

type Scheduler struct {
	DockerClient *client.Client

	Servers []*Server
}

func New() (*Scheduler, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	sched := new(Scheduler)

	sched.DockerClient = cli

	return sched, nil
}

func (s *Scheduler) ScheduleServer(ctx context.Context, input *ServerInput) (*Server, error) {
	_, err := s.DockerClient.ImagePull(ctx, "docker.io/library/openjdk:17", types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}

	hostDir := fmt.Sprintf("%s/servers/%s", viper.GetString("AARHUS_NODE_DATA"), input.Id)
	if err := utils.EnsureDirectory(hostDir); err != nil {
		return nil, err
	}

	cmd := []string{"java", fmt.Sprintf("-Xmx%dG", input.Ram), "-jar", input.ServerJar, "nogui"}

	resp, err := s.DockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image:      "openjdk:17",
			Cmd:        cmd,
			Tty:        false,
			WorkingDir: "/data",
			ExposedPorts: nat.PortSet{
				"25565/tcp": {},
			},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: hostDir,
					Target: "/data",
				},
			},
			PortBindings: nat.PortMap{
				"25565/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "25565", // Specify the host port to bind
					},
				},
			},
		},
		nil,
		nil,
		input.Id,
	)

	if err != nil {
		return nil, err
	}

	slog.Info("container created with id", slog.String("id", resp.ID))

	server := new(Server)

	server.Id = resp.ID

	s.Servers = append(s.Servers, server)

	return server, nil
}

func (s *Scheduler) StartAll(ctx context.Context) error {
	for _, server := range s.Servers {
		slog.Info("starting container", slog.String("id", server.Id))
		if err := s.DockerClient.ContainerStart(ctx, server.Id, types.ContainerStartOptions{}); err != nil {
			return err
		}
	}

	select {}

	return nil
}

func (s *Scheduler) Close() error {
	s.DockerClient.Close()
	return nil
}
