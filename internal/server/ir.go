package server

type IRPlayer struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Rank string `json:"rank"`
}
