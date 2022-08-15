package resource

import (
	"math"
	"math/rand"
	"time"
)

type Resource struct {
	name              string
	volumicMass       float64 // g/cm3
	averageItemVolume uint    // cm3
}

func (resource Resource) GetVolumicMass() float64 {
	return resource.volumicMass
}

func (resource Resource) GetName() string {
	return resource.name
}

func (resource *Resource) GenerateRandomItem() Item {
	source := rand.NewSource(time.Now().UnixNano())
	min := uint(math.Ceil(float64(resource.averageItemVolume) * 0.8))
	max := uint(math.Ceil(float64(resource.averageItemVolume) * 1.2))
	volume := uint((rand.New(source).NormFloat64()*0.2 + 1) * float64(resource.averageItemVolume))
	if volume > max {
		volume = max
	}
	if volume < min {
		volume = min
	}
	return Item{
		resource: resource,
		volume:   volume,
	}
}

var Apple = Resource{
	name:              "Apple",
	volumicMass:       0.96,
	averageItemVolume: 187,
}
