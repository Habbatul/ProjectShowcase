package controller

import (
	"awesomeProject/main/service"
	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	CategoryService *service.CategoryService
}

func NewCategoryController(categoryService *service.CategoryService) *CategoryController {
	return &CategoryController{CategoryService: categoryService}
}

// showAllCategoryNames godoc
// @Summary Get list category names
// @Description Retrieve all category names
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} object{categories=[]string}
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /category/names [get]
func (c *CategoryController) ShowAllCategoryNames(ctx *fiber.Ctx) error {
	categories, err := c.CategoryService.ShowAllCategoryNames()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something Wrong Happen"})
	}

	return ctx.JSON(fiber.Map{"categories": categories})
}
