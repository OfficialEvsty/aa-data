package usecase

type RaidItemDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Quantity int    `json:"quantity"`
}
