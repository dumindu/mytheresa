package product

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	e "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/model"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/mytheresa/go-hiring-challenge/internal/util/ctxutil"
	l "github.com/mytheresa/go-hiring-challenge/internal/util/logger"
)

type (
	API struct {
		repo   *repository.ProductsRepository
		logger *l.Logger
	}

	ListResponse struct {
		Products []*model.ProductResponse `json:"products"`
	}
)

func New(db *gorm.DB, logger *l.Logger) *API {
	return &API{
		repo:   repository.NewProductsRepository(db),
		logger: logger,
	}
}

func (api *API) GetAll(w http.ResponseWriter, r *http.Request) {
	reqID := ctxutil.RequestID(r.Context())

	products, err := api.repo.GetAllProducts()
	if err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespRepoDataAccessErr)
		return
	}

	resp := &ListResponse{
		Products: products.ToResponse(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeErr)
		return
	}
}
