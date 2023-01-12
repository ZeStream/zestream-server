package types

type Video struct {
	ID        string    `json:"id"`
	Src       string    `json:"src"`
	Type      string    `json:"type"`
	Watermark WaterMark `json:"watermark"`
}

type WaterMark struct {
	ID        string    `json:"id"`
	Src       string    `json:"src"`
	Dimension Dimension `json:"dimension"`
	Position  Dimension `json:"position"`
}

type Dimension struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type DIMENSION int

const (
	X DIMENSION = iota
	Y
)

var WaterMarkSizeMap = map[DIMENSION]string{
	X: "x",
	Y: "y",
}

var WaterMarkPositionMap = map[DIMENSION]string{
	X: "x",
	Y: "y",
}

func (w *WaterMark) IsEmpty() bool {
	if w.ID == "" {
		return true
	}

	return false
}
