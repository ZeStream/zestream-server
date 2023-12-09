package types

type Media struct {
	ID   string `json:"id" binding:"required"`
	Src  string `json:"src" binding:"required,url"`
	Type string `json:"type" binding:"required"`
}

type Video struct {
	Media
	Watermark WaterMark `json:"watermark"`
}

type Audio struct {
	Media
}

type WaterMark struct {
	ID        string    `json:"id" binding:"required_if=Watermark 1"`
	Src       string    `json:"src" binding:"required_if=Watermark 1"`
	Type      string    `json:"type" binding:"required_if=Watermark 1"`
	Dimension Dimension `json:"dimension" binding:"required_if=Watermark 1"`
	Position  Dimension `json:"position" binding:"required_if=Watermark 1"`
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
	return w.ID == ""
}
