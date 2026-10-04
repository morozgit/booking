package dto

type HotelAdd struct {
	Title    string `json:"title" validate:"required,min=2,max=100"`
	Location string `json:"location" validate:"required,min=2,max=400"`
}

type Hotel struct {
	Id       int `json:"id"`
	HotelAdd HotelAdd
}
