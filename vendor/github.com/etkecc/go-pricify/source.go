package pricify

import "encoding/json"

type sourceModel struct {
	Bases     []*sourceItem      `json:"bases"`
	Instances *sourceSectionItem `json:"instances"`
	Support   *sourceSectionItem `json:"support"`

	MatrixApps         []*sourceItem `json:"matrixApps"`
	MatrixBots         []*sourceItem `json:"matrixBots"`
	MatrixBridges      []*sourceItem `json:"matrixBridges"`
	MatrixBridgesVID   int64         `json:"matrixBridgesVID"`
	MatrixBridgesPrice int           `json:"matrixBridgesPrice"`
	MatrixAdditional   []*sourceItem `json:"additionalMatrixServices"`

	AdditionalServices []*sourceItem `json:"additionalServices"`
	ArchiveURL         string        `json:"archived"`
}

func (s *sourceModel) init() {
	if s.Instances == nil {
		s.Instances = &sourceSectionItem{}
	}
	if s.Support == nil {
		s.Support = &sourceSectionItem{}
	}
}

func (s *sourceModel) append(other *sourceModel) {
	if s == nil || other == nil {
		return
	}

	s.init()
	other.init()

	s.Bases = append(s.Bases, other.Bases...)
	s.Instances.Options = append(s.Instances.Options, other.Instances.Options...)
	s.Support.Options = append(s.Support.Options, other.Support.Options...)
	s.MatrixApps = append(s.MatrixApps, other.MatrixApps...)
	s.MatrixBots = append(s.MatrixBots, other.MatrixBots...)
	s.MatrixBridges = append(s.MatrixBridges, other.MatrixBridges...)
	if s.MatrixBridgesVID == 0 {
		s.MatrixBridgesVID = other.MatrixBridgesVID
	}
	if s.MatrixBridgesPrice == 0 {
		s.MatrixBridgesPrice = other.MatrixBridgesPrice
	}
	s.MatrixAdditional = append(s.MatrixAdditional, other.MatrixAdditional...)
	s.AdditionalServices = append(s.AdditionalServices, other.AdditionalServices...)
}

type sourceSectionItem struct {
	ID          string        `json:"id"`
	VID         int64         `json:"vid"`
	InventoryID string        `json:"iid"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Help        string        `json:"help"`
	Options     []*sourceItem `json:"options"`
	Archived    bool          `json:"archived"`
}

type sourceItem struct {
	ID          string   `json:"id"`          // Order form item ID
	VID         int64    `json:"vid"`         // V. ID
	InventoryID string   `json:"iid"`         // Inventory ID
	Name        string   `json:"name"`        // Human-readable name
	Description string   `json:"description"` // Human-readable description
	Help        string   `json:"help"`        // Help link (may not contain the full URL, just path)
	Price       int      `json:"price"`       // Price
	Regions     []string `json:"regions"`     // Regions list
	Archived    bool     `json:"archived"`    // Whether the item is archived
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

	data.fromSourceItem(source.Bases, "bases", "", "", "", 0, 0)
	data.fromSourceSection(source.Instances, "instances", 0)
	data.fromSourceSection(source.Support, "support", 0)

	data.fromSourceItem(source.MatrixApps, "matrix_apps", "Matrix Apps", "", "", 0, 0)
	data.fromSourceItem(source.MatrixBots, "matrix_bots", "Matrix Bots", "", "", 0, 0)
	data.fromSourceItem(source.MatrixBridges, "matrix_bridges", "Bridges", "With the help of bridges, you can access different networks right from your own Matrix server", "/help/bridges/", source.MatrixBridgesVID, source.MatrixBridgesPrice)
	data.fromSourceItem(source.MatrixAdditional, "matrix_additional", "Matrix Extras", "", "", 0, 0)

	data.fromSourceItem(source.AdditionalServices, "additional", "Extras", "", "", 0, 0)

	setCache(data)
	return data
}
