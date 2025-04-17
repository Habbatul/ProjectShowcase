package service

import (
	"awesomeProject/main/entity"
	"awesomeProject/main/model"
	"awesomeProject/main/repository"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ProjectService struct {
	ProjectRepository *repository.ProjectRepository
}

func NewProjectService(projectRepository *repository.ProjectRepository) *ProjectService {
	return &ProjectService{ProjectRepository: projectRepository}
}

func (s *ProjectService) GetProjectDetails(id uint) (model.ProjectResponse, error) {
	project, err := s.ProjectRepository.FindProjectByID(id)
	if err != nil {
		return model.ProjectResponse{}, err
	}

	var tags []string
	for _, tag := range project.Tags {
		tags = append(tags, tag.NameTag)
	}

	var categories []string
	for _, category := range project.Categories {
		categories = append(categories, category.Name)
	}

	var images []string
	for _, image := range project.Images {
		images = append(images, image.URLImg)
	}

	return model.ProjectResponse{
		Id:          project.ID,
		Name:        project.Name,
		Overview:    project.Overview,
		Description: project.Description,
		Note:        project.Note,
		URLProject:  project.URLProject,
		URLVideo:    project.URLVideo,
		OrderNumber: project.OrderNumber,
		Categories:  categories,
		Tags:        tags,
		Images:      images,
	}, nil
}

func (s *ProjectService) GetAllProject(cursor uint, categoryName string, limit int) ([]model.ProjectResponse, error) {

	projects, err := s.ProjectRepository.FindAllProject(limit, cursor, categoryName)
	if err != nil {
		return nil, err
	}

	var projectResponses []model.ProjectResponse

	for _, project := range projects {
		var tags []string
		for _, tag := range project.Tags {
			tags = append(tags, tag.NameTag)
		}

		var categories []string
		for _, catagory := range project.Categories {
			categories = append(categories, catagory.Name)
		}

		var images []string
		for _, image := range project.Images {
			images = append(images, image.URLImg)
		}

		projectResponses = append(projectResponses, model.ProjectResponse{
			Id:          project.ID,
			Name:        project.Name,
			Overview:    project.Overview,
			Description: project.Description,
			Note:        project.Note,
			URLProject:  project.URLProject,
			URLVideo:    project.URLVideo,
			OrderNumber: project.OrderNumber,
			Categories:  categories,
			Tags:        tags,
			Images:      images,
		})
	}

	return projectResponses, nil
}

func (s *ProjectService) CreateProject(dto *model.ProjectRequest) error {
	project := entity.Project{
		Name:        dto.Name,
		Overview:    dto.Overview,
		Description: dto.Description,
		Note:        dto.Note,
		URLProject:  dto.URLProject,
		URLVideo:    dto.URLVideo,
		OrderNumber: dto.OrderNumber,
	}

	for _, tagName := range dto.Tags {
		project.Tags = append(project.Tags, entity.Tag{NameTag: tagName})
	}

	for _, categoryName := range dto.Categories {
		project.Categories = append(project.Categories, entity.Category{Name: categoryName})
	}

	var baseUrl = "http://192.168.71.3:8080"
	for _, fileHeader := range dto.Images {
		fileName, err := saveImage(fileHeader, project.Name)
		if err != nil {
			return err
		}
		project.Images = append(project.Images, entity.Image{URLImg: baseUrl + "/project-images/" + fileName})
	}

	err := s.ProjectRepository.CreateProject(&project)
	return err
}

func saveImage(fileHeader *multipart.FileHeader, projectName string) (string, error) {
	uploadPath := "./uploads/images/"
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s_%s", projectName, fileHeader.Filename)
	filePath := filepath.Join(uploadPath, fileName)

	//open file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	//save to path
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = out.ReadFrom(file)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
