package entity

type Tag struct {
	ID       uint      `gorm:"primaryKey"`
	NameTag  string    `gorm:"unique" json:"name_tag"`
	Projects []Project `gorm:"many2many:project_tags;" json:"projects"`
}
