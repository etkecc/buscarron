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
	Value        string
	Section      string
	Price        int
	SectionPrice int
}

func (d *Data) fromSourceItem(sItems []sourceItem, section string, sectionPrice int) {
	for _, sItem := range sItems {
		item := &Item{
			ID:           sItem.ID,
			InventoryID:  sItem.InventoryID,
			Value:        "yes",
			Price:        sItem.Price,
			Section:      section,
			SectionPrice: sectionPrice,
		}
		d.items = append(d.items, item)
		d.idmap[item.ID] = item
		d.iidmap[item.InventoryID] = item
	}
}

func (d *Data) fromSourceSection(ssItem sourceSectionItem, section string, sectionPrice int) {
	for _, sItem := range ssItem.Options {
		item := &Item{
			ID:           ssItem.ID,
			InventoryID:  ssItem.InventoryID,
			Value:        sItem.ID,
			Price:        sItem.Price,
			Section:      section,
			SectionPrice: sectionPrice,
		}
		d.items = append(d.items, item)
		d.idmap[item.ID+item.Value] = item
		d.iidmap[item.InventoryID+item.Value] = item
	}
}

func (d *Data) find(key, value string) *Item {
	if item := d.idmap[key]; item != nil {
		return item
	}
	if item := d.iidmap[key]; item != nil {
		return item
	}
	if item := d.idmap[key+value]; item != nil {
		return item
	}
	if item := d.iidmap[key+value]; item != nil {
		return item
	}

	return nil
}

// Calculate total price based on input
func (d *Data) Calculate(input map[string]string) int {
	var total int
	sectionPriceAdded := map[string]bool{}
	for entry, value := range input {
		entry = strings.TrimSpace(strings.ToLower(entry))
		value = strings.TrimSpace(strings.ToLower(value))
		item := d.find(entry, value)
		if item == nil {
			continue
		}

		total += item.Price
		if item.SectionPrice > 0 && !sectionPriceAdded[item.Section] {
			total += item.SectionPrice
			sectionPriceAdded[item.Section] = true
		}
	}

	return total
}
