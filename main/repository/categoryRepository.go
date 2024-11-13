package repository

import (
	"awesomeProject/main/config"
	"awesomeProject/main/entity"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) FindAllCategoryNames() ([]string, error) {
	var categoryNames []string
	err := config.DB.Model(&entity.Category{}).Pluck("name", &categoryNames).Error
	if err != nil {
		return nil, err
	}

	return categoryNames, nil
}
