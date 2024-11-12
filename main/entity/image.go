package entity

type Image struct {
	ID        uint   `gorm:"primaryKey"`
	ProjectID uint   `json:"project_id"`
	URLImg    string `json:"url_img"`
}
