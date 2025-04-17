package entity

type Project struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `json:"name"`
	Overview    string     `json:"overview"`
	Description string     `json:"description"`
	Note        string     `json:"note"`
	URLProject  string     `json:"url_project"`
	URLVideo    string     `json:"url_video"`
	OrderNumber int32      `json:"order_number" gorm:"index"`
	Categories  []Category `gorm:"many2many:project_category;" json:"categories"`
	Tags        []Tag      `gorm:"many2many:project_tag;" json:"tags"`
	Images      []Image    `gorm:"foreignKey:ProjectID" json:"images"`
}
