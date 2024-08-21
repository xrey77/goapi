package models

type Photo struct {
	Photo_title string `json:"phototitle"`
	Photo_image string `json:"photoimage"`
}

type TempPhotos struct {
	ID          int64  `json:"id"`
	Photo_title string `json:"photo_title"`
	Photo_image string `json:"photo_image"`
}
