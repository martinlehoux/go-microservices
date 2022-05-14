package account

type AccountStore interface {
	Clear()
	Save(account Account) error
	LoadForIdentifier(identifier string) (Account, error)
}
