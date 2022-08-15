package inventory

import "go-microservices/resources/resource"

type Storage interface {
	Store(item resource.Item) error
	WouldStore(item resource.Item) error
}
