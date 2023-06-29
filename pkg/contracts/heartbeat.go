package contracts

type ServerHealthStatus int64

const (
	ServerHealthStatusUnknown ServerHealthStatus = iota
	ServerHealthStatusCrashed
	ServerHealthStatusRunning
	ServerHealthStatusStopped
)

type ServerHealthReport struct {
	ServerId string             `json:"server_id"`
	Status   ServerHealthStatus `json:"status"`

	// TODO: Add resource reports here
}

type HeartbeatInput struct {
	Networking CommonNetworking `json:"networking"`

	HealthReports []ServerHealthReport `json:"health_reports"`
	// TODO: Add resource reports here
}

type HeartbeatResponse struct {
	Tasks []*Task `json:"tasks"`
}
