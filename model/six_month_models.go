package model

type Info struct {
	Month string                     `json:"month"`
	Total int                        `json:"total"`
	Data  []Activeresignationforyear `json:"data"`
}
type Activeresignationforyear struct {
	Week    int `json:"week"`
	CountNo int `json:"count"`
}
