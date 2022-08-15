package inventory

import (
	"fmt"
	"math/rand"

	"go-microservices/resources/resource"
)

type Ore struct {
	Hardness           string
	MeltingTemperature int
	VolumicMass        float64 // g/cm3
}

var Cassiterite = Ore{
	Hardness:           "6-7",
	MeltingTemperature: 1200,
	VolumicMass:        7.0,
}

type OreBite struct {
	ore  *Ore
	Size float64 // cm3
}

// GetWeight returns the weight of the ore bite in grams
func GetWeight(bite OreBite) float64 {
	return bite.Size * bite.ore.VolumicMass
}

// Ore + Coal => Metal with T temp

func HarvestOre() (OreBite, bool) {
	if rand.Float32() < 0.9 {
		return OreBite{}, false
	}
	return OreBite{
		ore:  &Cassiterite,
		Size: 10,
	}, true
}

type Inventory struct {
	maxVolume uint // cm3
	items     []resource.Item
	owner     OwnerID
}

func NewInventory(owner OwnerID, maxVolume uint) Inventory {
	return Inventory{
		maxVolume: maxVolume,
		owner:     owner,
	}
}

type OwnerID string

func (inventory *Inventory) GetCurrentWeigth() uint {
	var weight uint
	for _, item := range inventory.items {
		weight += item.GetWeight()
	}
	return weight
}

func (inventory *Inventory) GetCurrentVolume() uint {
	var volume uint
	for _, item := range inventory.items {
		volume += item.GetVolume()
	}
	return volume
}

func (inventory *Inventory) WouldStore(item resource.Item) error {
	if inventory.maxVolume < inventory.GetCurrentVolume()+item.GetVolume() {
		return fmt.Errorf("inventory cannot store item %s of volume %d, because current volume is %d/%d", item.GetResource().GetName(), item.GetVolume(), inventory.GetCurrentVolume(), inventory.maxVolume)
	}
	return nil
}

func (inventory *Inventory) Store(item resource.Item) error {
	err := inventory.WouldStore(item)
	if err != nil {
		return err
	}
	inventory.items = append(inventory.items, item)
	return nil
}
