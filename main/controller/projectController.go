package controller

import (
	"awesomeProject/main/model"
	"awesomeProject/main/service"
	"github.com/gofiber/fiber/v2"
)

type ProjectController struct {
	ProjectService *service.ProjectService
}

func NewProjectController(projectService *service.ProjectService) *ProjectController {
	return &ProjectController{ProjectService: projectService}
}

// GetProjectDetails godoc
// @Summary Get project by ID
// @Description Retrieve project details by project ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} model.ProjectResponse
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /project/{id} [get]
func (pc *ProjectController) GetProjectDetails(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	project, err := pc.ProjectService.GetProjectDetails(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project not found"})
	}

	return ctx.JSON(project)
}

// GetAllProject godoc
// @Summary Get all projects with pagination and category filter
// @Description Retrieve project details with cursor-based pagination and optional category filter
// @Tags Projects
// @Accept json
// @Produce json
// @Param category query string false "Project Category Name"
// @Param cursor query int false "Cursor for pagination, default : 0"
// @Param cursor query int false "Limit for pagination, default : 6"
// @Success 200 {object} object{projects=[]model.ProjectResponse}
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /project [get]
func (pc *ProjectController) GetAllProject(ctx *fiber.Ctx) error {
	category := ctx.Query("category", "")

	cursor := ctx.QueryInt("cursor", 0)
	limit := ctx.QueryInt("limit", 6)

	projects, err := pc.ProjectService.GetAllProject(uint(cursor), category, limit)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Projects not found"})
	}

	return ctx.JSON(fiber.Map{"projects": projects})
}

// AddProject godoc
// @Summary Create a new project with images and tags
// @Description Create a new project with associated tags and images using form-data
// @Tags Projects
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Project Name"
// @Param overview formData string true "Project Overview"
// @Param description formData string true "Project Description"
// @Param note formData string false "Project Note"
// @Param url_project formData string false "Project URL"
// @Param categories[] formData []string false "Category item" collectionFormat(multi)
// @Param tags[] formData []string false "Tag item" collectionFormat(multi)
// @Param images[] formData []file false "Images" collectionFormat(multi)
// @Success 201 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /project [post]
func (pc *ProjectController) AddProject(ctx *fiber.Ctx) error {
	var projectRequest model.ProjectRequest

	projectRequest.Name = ctx.FormValue("name")
	projectRequest.Overview = ctx.FormValue("overview")
	projectRequest.Description = ctx.FormValue("description")
	projectRequest.Note = ctx.FormValue("note")
	projectRequest.URLProject = ctx.FormValue("url_project")

	//ambil semua multiple input collection array
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form data"})
	}

	//ambil multi tags
	tags := form.Value["tags[]"]
	if len(tags) > 0 {
		projectRequest.Tags = tags
	}

	//ambil multi categories
	categories := form.Value["categories[]"]
	if len(categories) > 0 {
		projectRequest.Categories = categories
	}

	//ambil multi images
	projectRequest.Images = form.File["images[]"]

	err = pc.ProjectService.CreateProject(&projectRequest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create project"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Project created successfully"})
}
