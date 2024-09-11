package pricify

import (
	"fmt"
	"strings"
)

// forbiddenValues is to help ignoring disabled items
var forbiddenValues = map[string]struct{}{
	"no":    {},
	"false": {},
}

// Data parsed from the source
type Data struct {
	items  []*Item
	idmap  map[string]*Item
	iidmap map[string]*Item
}

// Item is specific item parsed from the source
type Item struct {
	ID                 string
	InventoryID        string
	Name               string
	Description        string
	Help               string
	Value              string
	Price              int
	SectionID          string
	SectionName        string
	SectionDescription string
	SectionHelp        string
	SectionPrice       int
}

// fromSourceItem converts source items into the []*Item and adds them to the Data
func (d *Data) fromSourceItem(sItems []sourceItem, sectionID, sectionName, sectionDescription, sectionHelp string, sectionPrice int) {
	for _, sItem := range sItems {
		item := &Item{
			ID:                 sItem.ID,
			InventoryID:        sItem.InventoryID,
			Name:               sItem.Name,
			Description:        sItem.Description,
			Help:               sItem.Help,
			Value:              "yes",
			Price:              sItem.Price,
			SectionID:          sectionID,
			SectionName:        sectionName,
			SectionDescription: sectionDescription,
			SectionHelp:        sectionHelp,
			SectionPrice:       sectionPrice,
		}
		d.items = append(d.items, item)
		d.idmap[item.ID] = item
		d.iidmap[item.InventoryID] = item
	}
}

// fromSourceSection coverts source sections into the []*Item and adds them to the Data
func (d *Data) fromSourceSection(ssItem sourceSectionItem, sectionID string, sectionPrice int) {
	for _, sItem := range ssItem.Options {
		item := &Item{
			ID:           ssItem.ID,
			InventoryID:  ssItem.InventoryID,
			Value:        sItem.ID,
			Price:        sItem.Price,
			Name:         fmt.Sprintf("%s (%s)", ssItem.Name, sItem.Name),
			Description:  ssItem.Description,
			Help:         ssItem.Help,
			SectionID:    sectionID,
			SectionPrice: sectionPrice,
		}
		d.items = append(d.items, item)
		d.idmap[item.ID+item.Value] = item
		d.iidmap[item.InventoryID+item.Value] = item
	}
}

// find an item by key and value using IDs and Inventory IDs
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

// Calculate total price based on the input
func (d *Data) Calculate(input map[string]string) int {
	total, _ := d.CalculateVerbose(input)
	return total
}

// CalculateVerbose calculates total price and provides details about each found item
func (d *Data) CalculateVerbose(input map[string]string) (int, map[string]*Item) {
	var total int
	verbose := map[string]*Item{}

	// default value for etke_base_matrix
	if _, ok := input["etke_base_matrix"]; !ok {
		input["etke_base_matrix"] = "yes"
	}

	var withEmail bool
	for entry, value := range input {
		entry = strings.TrimSpace(strings.ToLower(entry))
		value = strings.TrimSpace(strings.ToLower(value))
		if _, ok := forbiddenValues[value]; ok {
			continue
		}
		item := d.find(entry, value)
		if item == nil {
			continue
		}
		if item.InventoryID == "etke_service_email" {
			withEmail = true
		}
	}

	sectionPriceAdded := map[string]bool{}
	for entry, value := range input {
		entry = strings.TrimSpace(strings.ToLower(entry))
		value = strings.TrimSpace(strings.ToLower(value))
		if _, ok := forbiddenValues[value]; ok {
			continue
		}

		item := d.find(entry, value)
		if item == nil {
			continue
		}

		// smtp relay should be free if email service is selected
		if item.InventoryID == "exim_relay_relay_use" && withEmail {
			freeRelay := *item
			freeRelay.Name = freeRelay.Name + " (free with email service)"
			freeRelay.Price = 0
			freeRelay.SectionPrice = 0

			item = &freeRelay
			value = "with email service"
		}

		if item.SectionPrice > 0 && !sectionPriceAdded[item.SectionID] {
			total += item.SectionPrice
			sectionPriceAdded[item.SectionID] = true
			verbose[item.SectionID] = &Item{
				ID:          "section-" + item.SectionID,
				InventoryID: "section_" + item.SectionID,
				Name:        item.SectionName,
				Description: item.SectionDescription,
				Help:        item.SectionHelp,
				Price:       item.SectionPrice,
			}
			continue
		}

		total += item.Price
		dup := *item
		dup.Value = value
		verbose[item.InventoryID] = &dup
	}

	return total, verbose
}
