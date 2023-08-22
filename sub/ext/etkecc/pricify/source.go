package pricify

import "encoding/json"

type sourceModel struct {
	Instances []sourceItem `json:"instances"`

	MatrixApps         []sourceItem `json:"matrixApps"`
	MatrixBots         []sourceItem `json:"matrixBots"`
	MatrixBridges      []sourceItem `json:"matrixBridges"`
	MatrixBridgesPrice int          `json:"matrixBridgesPrice"`

	AdditionalServices []sourceItem `json:"additionalServices"`
	AdvancedServices   []sourceItem `json:"advancedServices"`
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

	data.fromSource(source.Instances, "instances", 0)

	data.fromSource(source.MatrixApps, "matrix_apps", 0)
	data.fromSource(source.MatrixBots, "matrix_bots", 0)
	data.fromSource(source.MatrixBridges, "matrix_bridges", source.MatrixBridgesPrice)

	data.fromSource(source.AdditionalServices, "additional", 0)
	data.fromSource(source.AdvancedServices, "advanced", 0)

	return data
}
