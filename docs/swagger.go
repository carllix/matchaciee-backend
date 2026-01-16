package docs

import "github.com/google/uuid"

// Swagger DTO types for API documentation

// Response wrappers
type SwaggerResponse struct {
	Success bool `json:"success" example:"true"`
}

type SwaggerSuccessResponse struct {
	Success bool `json:"success" example:"true"`
	Data    any  `json:"data"`
}

type SwaggerErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Error message"`
}

type SwaggerValidationErrorResponse struct {
	Success bool              `json:"success" example:"false"`
	Error   string            `json:"error" example:"Validation failed"`
	Details map[string]string `json:"details"`
}

// Auth DTOs
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
	FullName string `json:"full_name" example:"John Doe"`
	Phone    string `json:"phone,omitempty" example:"+6281234567890"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email    string    `json:"email" example:"user@example.com"`
	FullName string    `json:"full_name" example:"John Doe"`
	Phone    *string   `json:"phone,omitempty" example:"+6281234567890"`
	Role     string    `json:"role" example:"member"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	Token        string       `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

type AuthSuccessResponse struct {
	Success bool         `json:"success" example:"true"`
	Data    AuthResponse `json:"data"`
}

type MeResponse struct {
	User UserResponse `json:"user"`
}

type MeSuccessResponse struct {
	Success bool       `json:"success" example:"true"`
	Data    MeResponse `json:"data"`
}

// Category DTOs
type CreateCategoryRequest struct {
	Name         string  `json:"name" example:"Matcha Drinks"`
	Slug         string  `json:"slug,omitempty" example:"matcha-drinks"`
	Description  *string `json:"description,omitempty" example:"Delicious matcha beverages"`
	DisplayOrder int     `json:"display_order,omitempty" example:"1"`
	IsActive     *bool   `json:"is_active,omitempty" example:"true"`
	ImageURL     *string `json:"image_url,omitempty" example:"https://example.com/image.jpg"`
}

type UpdateCategoryRequest struct {
	Name         *string `json:"name,omitempty" example:"Matcha Drinks Updated"`
	Slug         *string `json:"slug,omitempty" example:"matcha-drinks-updated"`
	Description  *string `json:"description,omitempty" example:"Updated description"`
	ImageURL     *string `json:"image_url,omitempty" example:"https://example.com/image.jpg"`
	DisplayOrder *int    `json:"display_order,omitempty" example:"2"`
	IsActive     *bool   `json:"is_active,omitempty" example:"true"`
}

type CategoryResponse struct {
	ID           uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string    `json:"name" example:"Matcha Drinks"`
	Slug         string    `json:"slug" example:"matcha-drinks"`
	Description  *string   `json:"description,omitempty" example:"Delicious matcha beverages"`
	ImageURL     *string   `json:"image_url,omitempty" example:"https://example.com/image.jpg"`
	DisplayOrder int       `json:"display_order" example:"1"`
	IsActive     bool      `json:"is_active" example:"true"`
	CreatedAt    string    `json:"created_at" example:"2025-01-07T10:00:00Z"`
	UpdatedAt    string    `json:"updated_at" example:"2025-01-07T10:00:00Z"`
}

type CategorySuccessResponse struct {
	Success bool             `json:"success" example:"true"`
	Data    CategoryResponse `json:"data"`
}

type CategoriesListResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Count      int                `json:"count" example:"5"`
}

type CategoriesSuccessResponse struct {
	Success bool                   `json:"success" example:"true"`
	Data    CategoriesListResponse `json:"data"`
}

// Product DTOs
type CreateCustomizationRequest struct {
	CustomizationType string  `json:"customization_type" example:"sweetness"`
	OptionName        string  `json:"option_name" example:"Less Sugar"`
	PriceModifier     float64 `json:"price_modifier" example:"0"`
	DisplayOrder      int     `json:"display_order,omitempty" example:"1"`
}

type UpdateCustomizationRequest struct {
	CustomizationType *string  `json:"customization_type,omitempty" example:"sweetness"`
	OptionName        *string  `json:"option_name,omitempty" example:"Extra Sugar"`
	PriceModifier     *float64 `json:"price_modifier,omitempty" example:"5000"`
	DisplayOrder      *int     `json:"display_order,omitempty" example:"2"`
}

type CreateProductRequest struct {
	Name            string                       `json:"name" example:"Matcha Latte"`
	Slug            string                       `json:"slug,omitempty" example:"matcha-latte"`
	Description     *string                      `json:"description,omitempty" example:"Creamy matcha latte"`
	CategoryID      *uuid.UUID                   `json:"category_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	BasePrice       float64                      `json:"base_price" example:"35000"`
	PreparationTime *int                         `json:"preparation_time,omitempty" example:"5"`
	DisplayOrder    int                          `json:"display_order,omitempty" example:"1"`
	IsAvailable     *bool                        `json:"is_available,omitempty" example:"true"`
	IsCustomizable  *bool                        `json:"is_customizable,omitempty" example:"true"`
	ImageURL        *string                      `json:"image_url,omitempty" example:"https://example.com/matcha.jpg"`
	Customizations  []CreateCustomizationRequest `json:"customizations,omitempty"`
}

type UpdateProductRequest struct {
	Name            *string    `json:"name,omitempty" example:"Matcha Latte Premium"`
	Slug            *string    `json:"slug,omitempty" example:"matcha-latte-premium"`
	Description     *string    `json:"description,omitempty" example:"Premium matcha latte"`
	BasePrice       *float64   `json:"base_price,omitempty" example:"40000"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	ImageURL        *string    `json:"image_url,omitempty" example:"https://example.com/matcha.jpg"`
	IsAvailable     *bool      `json:"is_available,omitempty" example:"true"`
	IsCustomizable  *bool      `json:"is_customizable,omitempty" example:"true"`
	PreparationTime *int       `json:"preparation_time,omitempty" example:"7"`
	DisplayOrder    *int       `json:"display_order,omitempty" example:"2"`
}

type CustomizationResponse struct {
	ID                uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CustomizationType string    `json:"customization_type" example:"sweetness"`
	OptionName        string    `json:"option_name" example:"Less Sugar"`
	PriceModifier     float64   `json:"price_modifier" example:"0"`
	DisplayOrder      int       `json:"display_order" example:"1"`
	CreatedAt         string    `json:"created_at" example:"2025-01-07T10:00:00Z"`
}

type ProductResponse struct {
	ID              uuid.UUID               `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name            string                  `json:"name" example:"Matcha Latte"`
	Slug            string                  `json:"slug" example:"matcha-latte"`
	Description     *string                 `json:"description,omitempty" example:"Creamy matcha latte"`
	Category        *CategoryResponse       `json:"category,omitempty"`
	BasePrice       float64                 `json:"base_price" example:"35000"`
	PreparationTime int                     `json:"preparation_time" example:"5"`
	DisplayOrder    int                     `json:"display_order" example:"1"`
	IsAvailable     bool                    `json:"is_available" example:"true"`
	IsCustomizable  bool                    `json:"is_customizable" example:"true"`
	ImageURL        *string                 `json:"image_url,omitempty" example:"https://example.com/matcha.jpg"`
	DeletedAt       *string                 `json:"deleted_at,omitempty"`
	Customizations  []CustomizationResponse `json:"customizations,omitempty"`
	CreatedAt       string                  `json:"created_at" example:"2025-01-07T10:00:00Z"`
	UpdatedAt       string                  `json:"updated_at" example:"2025-01-07T10:00:00Z"`
}

type ProductSuccessResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    ProductResponse `json:"data"`
}

type ProductsListResponse struct {
	Products []ProductResponse `json:"products"`
	Count    int               `json:"count" example:"10"`
}

type ProductsSuccessResponse struct {
	Success bool                 `json:"success" example:"true"`
	Data    ProductsListResponse `json:"data"`
}

type CustomizationSuccessResponse struct {
	Success bool                  `json:"success" example:"true"`
	Data    CustomizationResponse `json:"data"`
}

// Order DTOs
type OrderItemCustomization struct {
	CustomizationID uuid.UUID `json:"customization_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OptionName      string    `json:"option_name" example:"Less Sugar"`
}

type CreateOrderItemRequest struct {
	ProductID      uuid.UUID                `json:"product_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Quantity       int                      `json:"quantity" example:"2"`
	Notes          *string                  `json:"notes,omitempty" example:"Extra ice please"`
	Customizations []OrderItemCustomization `json:"customizations,omitempty"`
}

type CreateOrderRequest struct {
	CustomerName string                   `json:"customer_name" example:"John Doe"`
	Notes        *string                  `json:"notes,omitempty" example:"Please call when ready"`
	Items        []CreateOrderItemRequest `json:"items"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" example:"preparing" enums:"pending,preparing,ready,completed,cancelled"`
}

type UserSummary struct {
	ID       uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FullName string    `json:"full_name" example:"John Doe"`
	Email    string    `json:"email" example:"john@example.com"`
}

type OrderItemResponse struct {
	ID             uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProductName    string    `json:"product_name" example:"Matcha Latte"`
	Quantity       int       `json:"quantity" example:"2"`
	UnitPrice      float64   `json:"unit_price" example:"35000"`
	Subtotal       float64   `json:"subtotal" example:"70000"`
	Customizations any       `json:"customizations,omitempty"`
	Notes          *string   `json:"notes,omitempty" example:"Extra ice"`
}

type OrderResponse struct {
	ID           uuid.UUID           `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderNumber  string              `json:"order_number" example:"MC-250107-001"`
	CustomerName string              `json:"customer_name" example:"John Doe"`
	Status       string              `json:"status" example:"pending"`
	OrderSource  string              `json:"order_source" example:"member"`
	Subtotal     float64             `json:"subtotal" example:"70000"`
	Tax          float64             `json:"tax" example:"7000"`
	Total        float64             `json:"total" example:"77000"`
	Notes        *string             `json:"notes,omitempty" example:"Please call when ready"`
	Items        []OrderItemResponse `json:"items"`
	User         *UserSummary        `json:"user,omitempty"`
	CreatedAt    string              `json:"created_at" example:"2025-01-07T10:00:00Z"`
	CompletedAt  *string             `json:"completed_at,omitempty"`
}

type OrderSuccessResponse struct {
	Success bool          `json:"success" example:"true"`
	Data    OrderResponse `json:"data"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64           `json:"total" example:"100"`
	Page   int             `json:"page" example:"1"`
	Limit  int             `json:"limit" example:"20"`
}

type OrdersSuccessResponse struct {
	Success bool              `json:"success" example:"true"`
	Data    OrderListResponse `json:"data"`
}

// Payment DTOs
type PaymentTokenResponse struct {
	PaymentID   uuid.UUID `json:"payment_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Token       string    `json:"token" example:"66e4fa55-fdac-4ef9-91b5-733b97d1b862"`
	RedirectURL string    `json:"redirect_url" example:"https://app.sandbox.midtrans.com/snap/v2/vtweb/..."`
}

type PaymentSuccessResponse struct {
	Success bool                 `json:"success" example:"true"`
	Data    PaymentTokenResponse `json:"data"`
}

type MidtransWebhookRequest struct {
	TransactionTime   string  `json:"transaction_time" example:"2025-01-07 10:00:00"`
	TransactionStatus string  `json:"transaction_status" example:"settlement"`
	TransactionID     string  `json:"transaction_id" example:"1234567890"`
	StatusMessage     string  `json:"status_message" example:"Success"`
	StatusCode        string  `json:"status_code" example:"200"`
	SignatureKey      string  `json:"signature_key" example:"..."`
	PaymentType       string  `json:"payment_type" example:"gopay"`
	OrderID           string  `json:"order_id" example:"MC-250107-001-1234567890"`
	MerchantID        string  `json:"merchant_id" example:"G123456789"`
	GrossAmount       string  `json:"gross_amount" example:"77000.00"`
	FraudStatus       string  `json:"fraud_status" example:"accept"`
	Currency          string  `json:"currency" example:"IDR"`
	SettlementTime    *string `json:"settlement_time,omitempty" example:"2025-01-07 10:05:00"`
}

type WebhookSuccessResponse struct {
	Status string `json:"status" example:"success"`
}

type WebhookErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Error message"`
}

// Generic message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

type MessageSuccessResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    MessageResponse `json:"data"`
}
