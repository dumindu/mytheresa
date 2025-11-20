package product

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	e "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/model"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
)

type (
	API struct {
		repo *repository.ProductsRepository
	}

	ListResponse struct {
		Products []*model.ProductResponse `json:"products"`
	}
)

func New(db *gorm.DB) *API {
	return &API{
		repo: repository.NewProductsRepository(db),
	}
}

func (api *API) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := api.repo.GetAllProducts()
	if err != nil {
		// TODO: log error
		e.ServerError(w, e.RespRepoDataAccessErr)
		return
	}

	resp := &ListResponse{
		Products: products.ToResponse(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: log error
		e.ServerError(w, e.RespJSONEncodeErr)
		return
	}
}
