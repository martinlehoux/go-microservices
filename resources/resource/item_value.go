package resource

type Item struct {
	resource *Resource
	volume   uint // cm3
}

func NewItem(resource *Resource, volume uint) Item {
	return Item{
		resource: resource,
		volume:   volume,
	}
}

func (item *Item) GetWeight() uint {
	return uint(item.resource.GetVolumicMass() * float64(item.volume))
}

func (item *Item) GetResource() Resource {
	return *item.resource
}

func (item *Item) GetVolume() uint {
	return item.volume
}
