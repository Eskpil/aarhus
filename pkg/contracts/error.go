package contracts

type Error struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
}
