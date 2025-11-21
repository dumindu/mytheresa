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
