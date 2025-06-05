package repository

import (
	"awesomeProject/main/config"
	"awesomeProject/main/entity"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
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

// UpdateProject
func (r *ProjectRepository) UpdateProject(project *entity.Project) error {
	// Mulai transaksi
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if project.ID == 0 {
		tx.Rollback()
		return errors.New("project ID required")
	}
	log.Println(project.Note)

	//note, urlProject, urlVideo boleh kosong
	if err := tx.
		Model(&entity.Project{}).
		Where("id = ?", project.ID).
		Omit("Tags", "Categories", "Images").
		Updates(project).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.
		Model(&entity.Project{}).
		Where("id = ?", project.ID).
		Omit("Tags", "Categories", "Images").
		UpdateColumns(map[string]interface{}{
			"note":        project.Note,
			"url_project": project.URLProject,
			"url_video":   project.URLVideo,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	//file img -> insert; list id lama -> keep; list id lama hilang -> hapus;
	var existingImageIDs []uint
	if err := tx.Model(&entity.Image{}).Where("project_id =?", project.ID).Pluck("id", &existingImageIDs).Error; err != nil {
		tx.Rollback()
		return err
	}
	log.Println(existingImageIDs)

	imageAfterUpdate := make(map[uint]bool)
	for _, image := range project.Images {
		if image.ID != 0 {
			log.Print(image.ID)
			imageAfterUpdate[image.ID] = true
		} else {
			image.ProjectID = project.ID
			//kenapa kalo langsung image error
			if err := tx.Create(&image).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	for _, existingImageID := range existingImageIDs {
		//cek yang ada kalo di existing ada hapus
		if !imageAfterUpdate[existingImageID] {
			log.Println(existingImageID)
			if err := tx.Delete(&entity.Image{}, existingImageID).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	//handle tah
	var existingTags []entity.Tag
	if err := tx.Model(&project).Association("Tags").Find(&existingTags); err != nil {
		tx.Rollback()
		return err
	}
	newTagMap := make(map[string]bool)
	for _, t := range project.Tags {
		newTagMap[t.NameTag] = true
	}

	//bandingkan tag lama yang tidak ada lagi di tag baru (hapus kalo ga ada)
	for _, oldTag := range existingTags {
		if !newTagMap[oldTag.NameTag] {
			var tag entity.Tag
			if err := tx.Where("name_tag = ?", oldTag.NameTag).First(&tag).Error; err != nil {
				tx.Rollback()
				return err
			}

			var exists bool
			err := tx.Table("project_tag").
				Where("tag_id = ? AND project_id != ?", tag.ID, project.ID).
				Select("1").Limit(1).Scan(&exists).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			if !exists {
				if err := tx.Model(&tag).Association("Projects").Clear(); err != nil {
					tx.Rollback()
					return err
				}
				if err := tx.Delete(&tag).Error; err != nil {
					tx.Rollback()
					return err
				}
				fmt.Println("Tag dihapus dari tabel karena tidak digunakan:", oldTag.NameTag)
			}
		}
	}

	//menambahkan tag yang kalo pada tabel tag udah ada cuma relasiin kalo belum ada tambahin
	for i, tag := range project.Tags {
		var existingTag entity.Tag
		if err := tx.Where("name_tag = ?", tag.NameTag).
			First(&existingTag).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				//buat tag baru
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
			//gunakan tag yang sudah ada
			project.Tags[i] = existingTag
		}

	}

	//replace relasi many-to-many
	if err := tx.Model(&entity.Project{ID: project.ID}).
		Association("Tags").
		Replace(project.Tags); err != nil {
		tx.Rollback()
		return err
	}

	//handle category (sama ae kek diatas uehe)
	var existingCategories []entity.Category
	err := tx.Model(&project).Association("Categories").Find(&existingCategories)
	if err != nil {
		tx.Rollback()
		return err
	}

	newCategories := make(map[string]bool)
	for _, category := range project.Categories {
		newCategories[category.Name] = true
	}

	for _, oldCategory := range existingCategories {
		if !newCategories[oldCategory.Name] {
			var category entity.Category

			err = tx.Where("name = ?", oldCategory.Name).First(&category).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			var exists bool
			err = tx.Table("project_category").Where("category_id = ? AND project_id != ?", category.ID, project.ID).
				Select("1").Limit(1).Scan(&exists).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			if !exists {
				if err = tx.Model(&category).Association("Projects").Clear(); err != nil {
					tx.Rollback()
					return err
				}

				if err = tx.Delete(&category).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	for i, cat := range project.Categories {
		var existingCat entity.Category
		if err := tx.Where("name = ?", cat.Name).
			First(&existingCat).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := tx.Create(&cat).Error; err != nil {
					tx.Rollback()
					return err
				}
				project.Categories[i].ID = cat.ID
			} else {
				tx.Rollback()
				return err
			}
		} else {
			project.Categories[i] = existingCat
		}
	}
	if err = tx.Model(&entity.Project{ID: project.ID}).
		Association("Categories").
		Replace(project.Categories); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *ProjectRepository) FindAllImageURLs(projectID uint) []string {
	var fileNames []string

	config.DB.Model(&entity.Image{}).Where("project_id = ?", projectID).Pluck("url_img", &fileNames)

	return fileNames
}
