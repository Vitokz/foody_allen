package entity

type Chat struct {
	ID     int64  `json:"id" bson:"_id"`
	State  string `json:"state" bson:"state"`
	UserID int64  `json:"user_id" bson:"user_id"`
}

func (c *Chat) CollectionName() string {
	return "chats"
}
