package pricify

import (
	"strings"
)

type Data struct {
	items  []*Item
	idmap  map[string]*Item
	iidmap map[string]*Item
}

type Item struct {
	ID           string
	InventoryID  string
	Section      string
	Price        int
	SectionPrice int
}

var allowedValues = map[string]bool{
	"yes":     true,
	"true":    true,
	"on":      true,
	"synapse": true, // matrix_homeserver_implementation
}

func (d *Data) fromSource(sItems []sourceItem, section string, sectionPrice int) {
	for _, sItem := range sItems {
		item := &Item{
			ID:           sItem.ID,
			InventoryID:  sItem.InventoryID,
			Price:        sItem.Price,
			Section:      section,
			SectionPrice: sectionPrice,
		}
		d.items = append(d.items, item)
		d.idmap[item.ID] = item
		d.iidmap[item.InventoryID] = item
	}
}

// Calculate total price based on input
func (d *Data) Calculate(input map[string]string) int {
	var total int
	sectionPriceAdded := map[string]bool{}
	for entry, value := range input {
		entry = strings.TrimSpace(strings.ToLower(entry))
		value = strings.TrimSpace(strings.ToLower(value))
		if !allowedValues[value] {
			continue
		}

		item, ok := d.idmap[entry]
		if !ok {
			item, ok = d.iidmap[entry]
			if !ok {
				continue
			}
		}
		total += item.Price
		if item.SectionPrice > 0 && !sectionPriceAdded[item.Section] {
			total += item.SectionPrice
			sectionPriceAdded[item.Section] = true
		}
	}

	return total
}
