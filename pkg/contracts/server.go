package contracts

import "time"

type ServerType int64

const (
	ServerTypeUnknown ServerType = iota
	ServerTypeMinecraft
)

type ResourceAllocation struct {
	// in gigabytes
	Ram uint64 `json:"ram" bson:"ram"`

	CPU float64 `json:"cpu" bson:"cpu"`
}

type MinecraftServer struct {
	ServerJar string `json:"server_jar"`
}

type Server struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`

	NodeRef    string             `json:"node_ref" bson:"node_ref"`
	Type       ServerType         `json:"type" bson:"type"`
	Allocation ResourceAllocation `json:"allocation" bson:"allocation"`

	LastStarted time.Time `json:"last_started" bson:"last_started"`

	PortForwards []PortForward `json:"port_forwards" bson:"port_forwards"`

	MinecraftServer `json:"minecraft_server,omitempty" bson:"minecraft_server,omitempty"`
}

type MinecraftServerCreateInput struct {
	ServerJar string `json:"server_jar"`
}

type ServerCreateInput struct {
	Name    string     `json:"name"`
	NodeRef string     `json:"node_ref"`
	Type    ServerType `json:"server_type"`

	PortForwards []PortForward `json:"port_forwards"`

	MinecraftServerCreateInput `json:"minecraft_server"`
	Allocation                 ResourceAllocation `json:"allocation"`
}

type InternalServerInput struct {
	Id   string `json:"id"`
	Ram  uint64 `json:"ram"`
	Name string `json:"name"`

	ServerJar string `json:"server_jar"`
}
