package product

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	e "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/model"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
	"github.com/mytheresa/go-hiring-challenge/internal/util/ctxutil"
	l "github.com/mytheresa/go-hiring-challenge/internal/util/logger"
	"github.com/mytheresa/go-hiring-challenge/internal/util/requestutil"
)

type (
	API struct {
		repo   *repository.ProductRepository
		logger *l.Logger
	}

	ListResponse struct {
		Products []*model.ProductResponse `json:"products"`
		Total    int64                    `json:"total"`
	}
)

func New(db *gorm.DB, logger *l.Logger) *API {
	return &API{
		repo:   repository.NewProductRepository(db),
		logger: logger,
	}
}

// GetByCode
//
//	@summary		Get a product by code
//	@description	Retrieve a single product by its unique product code.
//	@tags			products
//
//	@router			/products/{code} [GET]
//	@accept			json
//	@produce		json
//
//	@param			code	path		string					true	"Product code (e.g. PROD001)"
//
//	@success		200		{object}	model.ProductResponse	"Product found"
//	@failure		400		{object}	e.Error					"Invalid or missing product code"
//	@failure		404		{object}	e.Error					"Product not found"
//	@failure		500		{object}	e.Error					"Internal server error"
func (api *API) GetByCode(w http.ResponseWriter, r *http.Request) {
	reqID := ctxutil.RequestID(r.Context())

	code := chi.URLParam(r, "code")
	if code == "" {
		e.BadRequest(w, e.RespInvalidCode)
		return
	}

	product, err := api.repo.GetOneByCode(code)
	if err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespRepoDataAccessErr)
		return
	}

	if product == nil {
		e.NotFound(w, e.RespNotFoundErr)
		return
	}

	if err := json.NewEncoder(w).Encode(product.ToResponse()); err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeErr)
		return
	}
}

// GetAll
//
//	@summary		List products
//	@description	Retrieve a paginated list of products, optionally filtered by category and maximum price.
//	@tags			products
//
//	@router			/products [GET]
//	@accept			json
//	@produce		json
//
//	@param			limit		query		int				false	"Maximum number of products per request (pagination limit)"
//	@param			offset		query		int				false	"Number of products to skip (pagination offset)"
//	@param			category	query		string			false	"Filter by category name (e.g. Clothing, Shoes, Accessories)"
//	@param			price-lt	query		number			false	"Return products/variants with effective price less than this value"
//
//	@success		200			{object}	ListResponse	"List of products with total count"
//	@failure		500			{object}	e.Error			"Internal server error"
func (api *API) GetAll(w http.ResponseWriter, r *http.Request) {
	reqID := ctxutil.RequestID(r.Context())

	limit, offset := requestutil.ParseQueryParamLimitOffset(r)
	priceLessThan, _ := decimal.NewFromString(r.URL.Query().Get("price-lt"))

	filter := &model.ProductFilter{
		Category:      r.URL.Query().Get("category"),
		PriceLessThan: priceLessThan,
	}

	products, err := api.repo.GetAllWithFilter(limit, offset, filter)
	if err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespRepoDataAccessErr)
		return
	}

	productsCount, err := api.repo.CountAllWithFilter(filter)
	if err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespRepoDataAccessErr)
		return
	}

	resp := &ListResponse{
		Products: products.ToResponse(),
		Total:    productsCount,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeErr)
		return
	}
}
