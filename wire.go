//go:build wireinject
// +build wireinject

package main

import (
	"awesomeProject/main/controller"
	"awesomeProject/main/repository"
	"awesomeProject/main/service"
	"github.com/google/wire"
)

func InitializeProject() (*controller.ProjectController, error) {
	wire.Build(
		controller.NewProjectController,
		service.NewProjectService,
		repository.NewProjectRepository,
	)
	return nil, nil
}

func InitializeCategory() (*controller.CategoryController, error) {
	wire.Build(
		controller.NewCategoryController,
		service.NewCategoryService,
		repository.NewCategoryRepository,
	)
	return nil, nil
}
