package category

import (
	"encoding/json"
	"net/http"
	"strconv"

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
		repo   *repository.CategoryRepository
		logger *l.Logger
	}

	ListResponse struct {
		Categories []*model.CategoryResponse `json:"categories"`
		Total      int64                     `json:"total"`
	}
)

func New(db *gorm.DB, logger *l.Logger) *API {
	return &API{
		repo:   repository.NewCategoryRepository(db),
		logger: logger,
	}
}

// Create
//
//	@summary		Create a product category
//	@description	Create a new product category using the provided payload.
//	@tags			categories
//
//	@router			/category [POST]
//	@accept			json
//	@produce		json
//	@param			body	body	model.CategoryForm	true	"Category create form"
//
//	@success		201
//	@failure		409	{object}	e.Error	"Category with the same code already exists"
//	@failure		422	{object}	e.Error	"Validation failed"
//	@failure		500	{object}	e.Error	"Internal server error"
func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	reqID := ctxutil.RequestID(r.Context())

	form := &model.CategoryForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONDecodeErr)
		return
	}

	// TODO: validate
	// e.UnprocessableEntity

	newCategory := form.ToModel()
	category, err := api.repo.Create(newCategory)
	if err != nil {
		if e.IsDuplicateDBEntry(err.Error()) {
			e.Conflict(w, e.RespConflictErr)
			return
		}

		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespRepoDataInsertErr)
		return
	}

	api.logger.Info().Str(l.KeyReqID, reqID).Str("id", strconv.Itoa(int(category.ID))).Msg("new category created")
	w.WriteHeader(http.StatusCreated)
}

// GetAll
//
//	@summary		List product categories
//	@description	Retrieve a paginated list of product categories.
//	@tags			categories
//
//	@router			/categories [GET]
//	@accept			json
//	@produce		json
//	@param			limit	query		int				false	"Maximum number of categories to return (pagination limit)"
//	@param			offset	query		int				false	"Number of categories to skip (pagination offset)"
//
//	@success		200		{object}	ListResponse	"List of categories with total count"
//	@failure		500		{object}	e.Error			"Internal server error"
func (api *API) GetAll(w http.ResponseWriter, r *http.Request) {
	reqID := ctxutil.RequestID(r.Context())

	limit, offset := requestutil.ParseQueryParamLimitOffset(r)

	categories, err := api.repo.GetAll(limit, offset)
	if err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespRepoDataAccessErr)
		return
	}

	categoryCount, err := api.repo.CountAll()
	if err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespRepoDataAccessErr)
		return
	}

	resp := &ListResponse{
		Categories: categories.ToResponse(),
		Total:      categoryCount,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		api.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeErr)
		return
	}
}
