package dao

type Row struct {
	Item   string
	Lang   float64
	Long   float64
	Images []string
	URL    string
}

type Radius struct {
	MinLongitude float64
	MaxLongitude float64
	MinLatitude  float64
	MaxLatitude  float64
}

type DAO interface {
	GetItemsByTerm(string) ([]Row, error)
	GetItemsByLocation(float64, float64) ([]Row, error)
	GetItemCount(string, Radius) (int, error)
}
