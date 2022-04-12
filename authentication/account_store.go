package authentication

type AccountStore interface {
	Save(account Account) error
	LoadForIdentifier(identifier string) (Account, error)
}
