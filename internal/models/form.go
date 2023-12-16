package models

type Form struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Owner       int64  `json:"owner"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PictureURL  string `json:"picture_url"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}
