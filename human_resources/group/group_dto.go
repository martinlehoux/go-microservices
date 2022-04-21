package group

type GroupDto struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	MembersCount int    `json:"members_count"`
}
