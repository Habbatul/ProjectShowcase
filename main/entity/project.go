package entity

type Project struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `json:"name"`
	Overview    string     `json:"overview"`
	Description string     `json:"description"`
	Note        string     `json:"note"`
	URLProject  string     `json:"url_project"`
	Categories  []Category `gorm:"many2many:project_category;" json:"categories"`
	Tags        []Tag      `gorm:"many2many:project_tag;" json:"tags"`
	Images      []Image    `gorm:"foreignKey:ProjectID" json:"images"`
}
