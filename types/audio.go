package types

type Audio struct {
	ID   string `json:"id" binding:"required"`
	Src  string `json:"src" binding:"required,url"`
	Type string `json:"type" binding:"required"`
}
