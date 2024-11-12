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

func (r *ProjectRepository) GetProjectByID(id uint) (entity.Project, error) {
	var project entity.Project
	err := config.DB.Preload("Tags").Preload("Images").First(&project, id).Error
	if err != nil {
		return project, err
	}
	return project, nil
}

func (r *ProjectRepository) CreateProjectWithTagsAndImages(project *entity.Project) error {
	tx := config.DB.Begin()

	if err := tx.Omit("Tags").Create(&project).Error; err != nil {
		tx.Rollback()
		return err
	}

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

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
