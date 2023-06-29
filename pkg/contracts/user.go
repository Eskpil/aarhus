package contracts

type User struct {
	Id       string `json:"id" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Verified bool   `json:"verified" bson:"verified"`
	Avatar   string `json:"avatar" bson:"avatar"`
}
