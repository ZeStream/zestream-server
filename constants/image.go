package constants

type IMAGE_QUERY int

const (
	Width IMAGE_QUERY = iota
	Height
)

var ImageQuery = map[IMAGE_QUERY]string{
	Width:  "width",
	Height: "height",
}
