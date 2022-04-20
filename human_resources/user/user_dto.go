package user

type UserDto struct {
	ID            string `json:"id"`
	PreferredName string `json:"preferred_name"`
	Email         string `json:"email"`
}

func DtoFrom(user User) UserDto {
	return UserDto{
		ID:            user.id.String(),
		PreferredName: user.preferredName,
		Email:         user.email,
	}
}

type UserListDto struct {
	Items []UserDto `json:"items"`
	Total int       `json:"total"`
}
