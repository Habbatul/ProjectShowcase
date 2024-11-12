package entity

type Project struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Note        string  `json:"note"`
	URLProject  string  `json:"url_project"`
	Tags        []Tag   `gorm:"many2many:project_tags;" json:"tags"`
	Images      []Image `gorm:"foreignKey:ProjectID" json:"images"`
}
