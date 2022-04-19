package user

type UserStore interface {
	Load(userId UserID) (User, error)
	Save(user User) error
	GetMany() ([]User, error)
	GetByEmail(email string) (User, error)
	EmailExists(email string) (bool, error)
}
