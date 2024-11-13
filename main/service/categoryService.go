package service

import (
	"awesomeProject/main/repository"
)

type CategoryService struct {
	CategoryRepository *repository.CategoryRepository
}

func NewCategoryService(categoryRepository *repository.CategoryRepository) *CategoryService {
	return &CategoryService{CategoryRepository: categoryRepository}
}

func (s *CategoryService) ShowAllCategoryNames() ([]string, error) {
	var categories []string

	categories, err := s.CategoryRepository.FindAllCategoryNames()
	if err != nil {
		return categories, err
	}

	return categories, nil
}
