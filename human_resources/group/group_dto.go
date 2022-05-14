package group

type GroupDto struct {
	ID           string `json:"id" bson:"id"`
	Name         string `json:"name" bson:"name"`
	Description  string `json:"description" bson:"description"`
	MembersCount int    `json:"members_count" bson:"members_count"`
}
