package settlement

import (
	postgres "github.com/failuretoload/datamonster/settlement/internal"
	"github.com/failuretoload/datamonster/store"
	"github.com/failuretoload/datamonster/web"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	repo *postgres.PostgresRepo
}

func NewController(conn store.Connection) *Controller {
	repo := postgres.New(conn)
	return &Controller{repo: repo}
}

type SettlementDTO struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	SurvivalLimit       int    `json:"limit"`
	DepartingSurvival   int    `json:"departing"`
	CollectiveCognition int    `json:"cc"`
	Year                int    `json:"year"`
}

func withPermission(routeHandler http.HandlerFunc) http.HandlerFunc {
	return web.ValidatePermissions([]string{"manage:settlements"}, routeHandler)
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Get("/settlement", withPermission(c.getSettlements))
	r.Post("/settlement", withPermission(c.createSettlement))
	r.Route("/settlement/{id}", func(r chi.Router) {
		r.Get("/", withPermission(c.getSettlement))
	})
}

func (c Controller) getSettlements(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(web.UserIdKey).(string)
	settlements, repoErr := c.repo.Select(r.Context(), userID)
	if repoErr != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Error retrieving settlements")
		return
	}
	data := domainListToDto(settlements)
	web.MakeJsonResponse(w, http.StatusOK, data)
}

type CreateSettlementRequest struct {
	Name string `json:"name"`
}

func (c Controller) createSettlement(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(web.UserIdKey).(string)
	if !ok || userID == "" {
		web.MakeJsonResponse(w, http.StatusBadRequest, "no user id provided")
		return
	}
	var body CreateSettlementRequest
	err := web.DecodeJsonRequest(r.Body, &body)
	if err != nil {
		web.MakeJsonResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" {
		web.MakeJsonResponse(w, http.StatusBadRequest, "name is required")
		return
	}
	settlement := postgres.Settlement{
		Owner:               userID,
		Name:                body.Name,
		SurvivalLimit:       1,
		DepartingSurvival:   0,
		CollectiveCognition: 0,
		CurrentYear:         1,
	}
	newId, insertErr := c.repo.Insert(r.Context(), settlement)
	if insertErr != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Unable to create settlement")
		return
	}

	settlement.Id = newId
	dto := domainToDto(settlement)
	web.MakeJsonResponse(w, http.StatusOK, dto)
}

func (c Controller) getSettlement(w http.ResponseWriter, r *http.Request) {
	settlementId := chi.URLParam(r, "id")
	settlement, repoErr := c.repo.Get(r.Context(), settlementId)
	if repoErr != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Error retrieving settlement")
		return
	}
	dto := domainToDto(settlement)
	web.MakeJsonResponse(w, http.StatusOK, dto)
}

func domainListToDto(settlements []postgres.Settlement) []SettlementDTO {
	dtos := []SettlementDTO{}
	for _, s := range settlements {
		dto := SettlementDTO{
			Id:                  s.Id,
			Name:                s.Name,
			SurvivalLimit:       s.SurvivalLimit,
			DepartingSurvival:   s.DepartingSurvival,
			CollectiveCognition: s.CollectiveCognition,
			Year:                s.CurrentYear,
		}
		dtos = append(dtos, dto)
	}
	return dtos
}

func domainToDto(s postgres.Settlement) SettlementDTO {
	return SettlementDTO{
		Id:                  s.Id,
		Name:                s.Name,
		SurvivalLimit:       s.SurvivalLimit,
		DepartingSurvival:   s.DepartingSurvival,
		CollectiveCognition: s.CollectiveCognition,
		Year:                s.CurrentYear,
	}
}
