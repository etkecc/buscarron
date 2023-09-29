package pricify

import "encoding/json"

type sourceModel struct {
	Bases     []sourceItem      `json:"bases"`
	Instances sourceSectionItem `json:"instances"`
	Support   sourceSectionItem `json:"support"`

	MatrixApps         []sourceItem `json:"matrixApps"`
	MatrixBots         []sourceItem `json:"matrixBots"`
	MatrixBridges      []sourceItem `json:"matrixBridges"`
	MatrixBridgesPrice int          `json:"matrixBridgesPrice"`

	AdditionalServices []sourceItem `json:"additionalServices"`
	AdvancedServices   []sourceItem `json:"advancedServices"`
}

type sourceSectionItem struct {
	ID          string       `json:"id"`
	InventoryID string       `json:"iid"`
	Options     []sourceItem `json:"options"`
}

type sourceItem struct {
	ID          string `json:"id"`
	InventoryID string `json:"iid"`
	Price       int    `json:"price"`
}

func parseSource(data []byte) (*sourceModel, error) {
	var source *sourceModel
	err := json.Unmarshal(data, &source)
	return source, err
}

func convertToData(source *sourceModel) *Data {
	data := &Data{
		items:  []*Item{},
		idmap:  map[string]*Item{},
		iidmap: map[string]*Item{},
	}

	data.fromSourceItem(source.Bases, "bases", 0)
	data.fromSourceSection(source.Instances, "instances", 0)
	data.fromSourceSection(source.Support, "support", 0)

	data.fromSourceItem(source.MatrixApps, "matrix_apps", 0)
	data.fromSourceItem(source.MatrixBots, "matrix_bots", 0)
	data.fromSourceItem(source.MatrixBridges, "matrix_bridges", source.MatrixBridgesPrice)

	data.fromSourceItem(source.AdditionalServices, "additional", 0)
	data.fromSourceItem(source.AdvancedServices, "advanced", 0)

	return data
}
