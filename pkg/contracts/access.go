package contracts

type AccessEntry struct {
	Id    string `json:"id" bson:"_id"`
	Email string `json:"email" bson:"id"`
}
