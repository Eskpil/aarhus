package contracts

import "time"

type NodeTokenStatus int64

const (
	NodeTokenStatusRevoked NodeTokenStatus = iota
	NodeTokenStatusActive
)

type NodeToken struct {
	Value  string          `json:"value"`
	Status NodeTokenStatus `json:"status"`
}

type Node struct {
	Id string `json:"id" bson:"_id"`

	Networking CommonNetworking `json:"networking" bson:"networking"`

	Name     string `json:"name"`
	Location string `json:"location"`

	Observed      bool      `json:"observed"`
	LastHeartbeat time.Time `json:"last_heartbeat" bson:"last_heartbeat"`

	Tokens []NodeToken `json:"tokens"`

	Tasks []Task `json:"tasks"`
}

type NodeCreateInput struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}
