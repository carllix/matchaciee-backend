package repositories

import (
	"errors"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrProductSlugExists = errors.New("product slug already exists")
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindByUUID(uuid uuid.UUID) (*models.Product, error)
	FindByUUIDIncludingDeleted(uuid uuid.UUID) (*models.Product, error)
	FindBySlug(slug string) (*models.Product, error)
	FindAll(includeDeleted bool, isAvailable *bool, categoryID *uint) ([]models.Product, error)
	FindByCategoryUUID(categoryUUID uuid.UUID, includeDeleted bool, isAvailable *bool) ([]models.Product, error)
	Update(product *models.Product) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	ExistsBySlug(slug string) (bool, error)

	CreateCustomization(customization *models.ProductCustomization) error
	FindCustomizationByUUID(uuid uuid.UUID) (*models.ProductCustomization, error)
	UpdateCustomization(customization *models.ProductCustomization) error
	DeleteCustomization(id uint) error
	FindCustomizationsByProductID(productID uint) ([]models.ProductCustomization, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	exists, err := r.ExistsBySlug(product.Slug)
	if err != nil {
		return err
	}
	if exists {
		return ErrProductSlugExists
	}

	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Preload("Category").
		Preload("Customizations", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindByUUID(uuid uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Preload("Category").
		Preload("Customizations", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Where("uuid = ? AND deleted_at IS NULL", uuid).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindByUUIDIncludingDeleted(uuid uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Unscoped(). // Include soft-deleted records
		Preload("Category").
		Preload("Customizations", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Where("uuid = ?", uuid).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySlug(slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Preload("Category").
		Preload("Customizations", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Where("slug = ? AND deleted_at IS NULL", slug).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(includeDeleted bool, isAvailable *bool, categoryID *uint) ([]models.Product, error) {
	var products []models.Product
	query := r.db.
		Preload("Category").
		Preload("Customizations", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		})

	// Filter by soft delete
	if !includeDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	// Filter by availability
	if isAvailable != nil {
		query = query.Where("is_available = ?", *isAvailable)
	}

	// Filter by category
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	err := query.Order("display_order ASC, created_at DESC").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindByCategoryUUID(categoryUUID uuid.UUID, includeDeleted bool, isAvailable *bool) ([]models.Product, error) {
	// Find the category to get its internal ID
	var category models.Category
	if err := r.db.Where("uuid = ?", categoryUUID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return r.FindAll(includeDeleted, isAvailable, &category.ID)
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) SoftDelete(id uint) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *productRepository) Restore(id uint) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *productRepository) HardDelete(id uint) error {
	return r.db.Unscoped().Delete(&models.Product{}, id).Error
}

func (r *productRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("slug = ? AND deleted_at IS NULL", slug).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *productRepository) CreateCustomization(customization *models.ProductCustomization) error {
	return r.db.Create(customization).Error
}

func (r *productRepository) FindCustomizationByUUID(uuid uuid.UUID) (*models.ProductCustomization, error) {
	var customization models.ProductCustomization
	err := r.db.Where("uuid = ?", uuid).First(&customization).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customization not found")
		}
		return nil, err
	}
	return &customization, nil
}

func (r *productRepository) UpdateCustomization(customization *models.ProductCustomization) error {
	return r.db.Save(customization).Error
}

func (r *productRepository) DeleteCustomization(id uint) error {
	return r.db.Delete(&models.ProductCustomization{}, id).Error
}

func (r *productRepository) FindCustomizationsByProductID(productID uint) ([]models.ProductCustomization, error) {
	var customizations []models.ProductCustomization
	err := r.db.
		Where("product_id = ?", productID).
		Order("display_order ASC, created_at ASC").
		Find(&customizations).Error
	if err != nil {
		return nil, err
	}
	return customizations, nil
}
