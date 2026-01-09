package repositories

import (
	"errors"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategorySlugExists = errors.New("category slug already exists")
	ErrCategoryNameExists = errors.New("category name already exists")
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindByID(id uint) (*models.Category, error)
	FindByUUID(uuid uuid.UUID) (*models.Category, error)
	FindBySlug(slug string) (*models.Category, error)
	FindAll(isActive *bool) ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	ExistsBySlug(slug string) (bool, error)
	ExistsByName(name string) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	exists, err := r.ExistsBySlug(category.Slug)
	if err != nil {
		return err
	}
	if exists {
		return ErrCategorySlugExists
	}

	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindByUUID(uuid uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("uuid = ?", uuid).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindAll(isActive *bool) ([]models.Category, error) {
	var categories []models.Category
	query := r.db.Order("display_order ASC, created_at DESC")

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Model(&models.Category{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *categoryRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *categoryRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
