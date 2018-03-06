package dao

type Row struct {
	ItemID   string   `json:"item_id"`
	Lat      float64  `json:"latitude"`
	Long     float64  `json:"longitude"`
	URL      string   `json:"item_url"`
	Images   []string `json:"item_images"`
	Distance float64  `json:"-"`
}

type Rows []Row

func (r Rows) Len() int {
	return len(r)
}

func (r Rows) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Rows) Less(i, j int) bool {
	return r[i].Distance < r[j].Distance
}

type Radius struct {
	MinLatitude     float64
	MaxLatitude     float64
	MinLongitude    float64
	MaxLongitude    float64
	CenterLatitude  float64
	CenterLongitude float64
}

type DAO interface {
	GetItems(string, Radius) ([]Row, error)
	CountItems(string, Radius) (int, error)
}
