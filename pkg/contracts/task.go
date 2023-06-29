package contracts

type TaskStatus uint64
type TaskType uint64

const (
	TaskTypeUnknown TaskType = iota
	TaskTypeForwardPort
	TaskTypeCreateServer
)

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusScheduled
	TaskStatusFinished
	TaskStatusAborted
)

type Task struct {
	Id     string     `json:"id"`
	Type   TaskType   `json:"type"`
	Status TaskStatus `json:"status"`

	TaskForwardPort  `json:"task_forward_port"`
	TaskCreateServer `json:"task_create_server"`
}

type TaskForwardPort struct {
	ServerId    string `json:"server_id"`
	PortForward `json:"port_forward"`
}

type TaskCreateServer struct {
	ServerId string `json:"server_id"`
}

type UpdateTaskStatus struct {
	Status TaskStatus `json:"status"`
}
