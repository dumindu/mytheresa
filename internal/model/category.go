package model

func (c *Category) TableName() string {
	return "categories"
}

type (
	Categories []*Category
	Category   struct {
		ID   uint   `gorm:"primaryKey"`
		Code string `gorm:"uniqueIndex;not null"`
		Name string `gorm:"not null"`
	}

	CategoryResponse struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	CategoryForm struct {
		Code string `json:"code" validate:"required,max=32"`
		Name string `json:"name" validate:"required,max=256"`
	}
)

func (cs Categories) ToResponse() []*CategoryResponse {
	categories := make([]*CategoryResponse, len(cs))
	for i, c := range cs {
		categories[i] = c.ToResponse()
	}

	return categories
}

func (c *Category) ToResponse() *CategoryResponse {
	return &CategoryResponse{
		Code: c.Code,
		Name: c.Name,
	}
}

func (f *CategoryForm) ToModel() *Category {
	return &Category{
		Code: f.Code,
		Name: f.Name,
	}
}
