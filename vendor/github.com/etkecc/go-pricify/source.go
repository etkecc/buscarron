package pricify

import "encoding/json"

type sourceModel struct {
	Bases     []*sourceItem      `json:"bases"`
	Instances *sourceSectionItem `json:"instances"`
	Support   *sourceSectionItem `json:"support"`

	MatrixApps         []*sourceItem `json:"matrixApps"`
	MatrixBots         []*sourceItem `json:"matrixBots"`
	MatrixBridges      []*sourceItem `json:"matrixBridges"`
	MatrixBridgesPrice int           `json:"matrixBridgesPrice"`
	MatrixAdditional   []*sourceItem `json:"additionalMatrixServices"`

	AdditionalServices []*sourceItem `json:"additionalServices"`
}

type sourceSectionItem struct {
	ID          string       `json:"id"`
	InventoryID string       `json:"iid"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Help        string       `json:"help"`
	Options     []sourceItem `json:"options"`
}

type sourceItem struct {
	ID          string `json:"id"`          // Order form item ID
	InventoryID string `json:"iid"`         // Inventory ID
	Name        string `json:"name"`        // Human-readable name
	Description string `json:"description"` // Human-readable description
	Help        string `json:"help"`        // Help link (may not contain the full URL, just path)
	Price       int    `json:"price"`       // Price
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

	data.fromSourceItem(source.Bases, "bases", "", "", "", 0)
	data.fromSourceSection(source.Instances, "instances", 0)
	data.fromSourceSection(source.Support, "support", 0)

	data.fromSourceItem(source.MatrixApps, "matrix_apps", "", "", "", 0)
	data.fromSourceItem(source.MatrixBots, "matrix_bots", "", "", "", 0)
	data.fromSourceItem(source.MatrixBridges, "matrix_bridges", "Bridges", "With the help of bridges, you can access different networks right from your own Matrix server", "/help/bridges/", source.MatrixBridgesPrice)
	data.fromSourceItem(source.MatrixAdditional, "matrix_additional", "", "", "", 0)

	data.fromSourceItem(source.AdditionalServices, "additional", "", "", "", 0)

	return data
}
