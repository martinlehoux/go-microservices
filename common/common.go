package common

import "github.com/google/uuid"

type ID = uuid.UUID

func CreateID() ID {
	return uuid.New()
}
