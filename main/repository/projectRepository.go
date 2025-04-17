package repository

import (
	"awesomeProject/main/config"
	"awesomeProject/main/entity"
	"errors"
	"gorm.io/gorm"
)

type ProjectRepository struct{}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{}
}

func (r *ProjectRepository) FindProjectByID(id uint) (entity.Project, error) {
	var project entity.Project
	err := config.DB.Preload("Tags").Preload("Categories").Preload("Images").First(&project, id).Error
	if err != nil {
		return project, err
	}
	return project, nil
}

func (r *ProjectRepository) FindAllProject(limit int, cursor uint, categoryName string) ([]entity.Project, error) {
	var projects []entity.Project
	query := config.DB.Preload("Tags").Preload("Categories").Preload("Images")

	if categoryName != "" {
		query = query.Joins("JOIN project_category ON project_category.project_id = projects.id").
			Joins("JOIN categories ON categories.id = project_category.category_id").
			Where("categories.name = ?", categoryName)
	}

	if cursor > 0 {
		query = query.Where("projects.order_number < ?", cursor)
	}

	err := query.Order("projects.order_number DESC").Limit(limit).Find(&projects).Error
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) CreateProject(project *entity.Project) error {
	tx := config.DB.Begin()

	//kalo kosong isi max order_number + 1
	if project.OrderNumber == 0 {
		var maxOrderNumber int32
		err := tx.Model(&entity.Project{}).Select("COALESCE(MAX(order_number), 0)").Scan(&maxOrderNumber).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		project.OrderNumber = maxOrderNumber + 1
	}

	//simpan project dan image
	if err := tx.Omit("Tags", "Categories").Create(&project).Error; err != nil {
		tx.Rollback()
		return err
	}

	//simpan tags
	for i, tag := range project.Tags {
		var existingTag entity.Tag

		if err := tx.Where("name_tag = ?", tag.NameTag).First(&existingTag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := tx.Create(&tag).Error; err != nil {
					tx.Rollback()
					return err
				}

				project.Tags[i].ID = tag.ID
			} else {

				tx.Rollback()
				return err
			}
		} else {
			project.Tags[i] = existingTag
		}
	}

	if err := tx.Model(&project).Association("Tags").Replace(project.Tags); err != nil {
		tx.Rollback()
		return err
	}

	//simpan category
	for i, category := range project.Categories {
		var existingCategory entity.Category

		if err := tx.Where("name = ?", category.Name).First(&existingCategory).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := tx.Create(&category).Error; err != nil {
					tx.Rollback()
					return err
				}

				project.Categories[i].ID = category.ID
			} else {
				tx.Rollback()
				return err
			}
		} else {
			project.Categories[i] = existingCategory
		}
	}

	if err := tx.Model(&project).Association("Categories").Replace(project.Categories); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
