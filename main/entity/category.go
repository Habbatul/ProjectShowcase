package entity

type Category struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"unique" json:"category"`
	Projects []Project `gorm:"many2many:project_category;" json:"projects"`
}
