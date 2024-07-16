package survivor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/failuretoload/datamonster/store"
	repo "github.com/failuretoload/datamonster/survivor/internal"
	"github.com/failuretoload/datamonster/web"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	db *repo.PostGresRepo
}

func NewController(conn store.Connection) *Controller {
	r := repo.NewRepo(conn)
	return &Controller{db: r}
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Use(settlementIdExtractor)
	r.Get("/settlement/{id}/survivor", withPermission(c.getSurvivors))
	r.Post("/settlement/{id}/survivor", withPermission(c.createSurvivor))
}

func (c Controller) getSurvivors(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	settlementId, convErr := strconv.Atoi(param)
	if convErr != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "unable to convert query param")
		return
	}
	survivors, err := c.db.GetAllSurvivorsForSettlement(r.Context(), settlementId)
	if err != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "Error retrieving survivors")
		return
	}
	data := dtoListFromDomain(survivors)
	web.MakeJsonResponse(w, http.StatusOK, data)
}

func (c Controller) createSurvivor(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	settlementId, convErr := strconv.Atoi(param)
	if convErr != nil {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "settlement id must be a positive integer")
		return
	}
	survivorDTO := SurvivorDTO{}
	decodeErr := json.NewDecoder(r.Body).Decode(&survivorDTO)
	if decodeErr != nil || survivorDTO.Name == "" {
		web.MakeJsonResponse(w, http.StatusInternalServerError, "unable to decode request body")
		return
	}
	survivorDTO.Settlement = settlementId
	err := c.db.CreateSurvivor(r.Context(), domainFromDTO(survivorDTO))
	if err != nil {
		dupError := repo.DuplicateNameError{}
		if errors.As(err, &dupError) {
			web.MakeJsonResponse(w, http.StatusBadRequest, fmt.Sprintf("survivor with name %s already exists", survivorDTO.Name))
			return
		}
		web.MakeJsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("error creating survivor %s", survivorDTO.Name))
		return
	}
	web.MakeJsonResponse(w, http.StatusNoContent, nil)
}

func withPermission(routeHandler http.HandlerFunc) http.HandlerFunc {
	return web.ValidatePermissions([]string{"manage:survivors"}, routeHandler)
}

type SurvivorDTO struct {
	Id               int     `json:"id"`
	Settlement       int     `json:"settlement"`
	Name             string  `json:"name"`
	Birth            int     `json:"birth"`
	Gender           string  `json:"gender"`
	Status           *string `json:"status,omitempty"`
	HuntXp           int     `json:"huntXp"`
	Survival         int     `json:"survival"`
	Movement         int     `json:"movement"`
	Accuracy         int     `json:"accuracy"`
	Strength         int     `json:"strength"`
	Evasion          int     `json:"evasion"`
	Luck             int     `json:"luck"`
	Speed            int     `json:"speed"`
	Insanity         int     `json:"insanity"`
	SystemicPressure int     `json:"systemicPressure"`
	Torment          int     `json:"torment"`
	Lumi             int     `json:"lumi"`
	Courage          int     `json:"courage"`
	Understanding    int     `json:"understanding"`
}

func dtoFromDomain(s repo.Survivor) SurvivorDTO {
	return SurvivorDTO(s)
}

func domainFromDTO(s SurvivorDTO) repo.Survivor {
	return repo.Survivor(s)
}

func dtoListFromDomain(s []repo.Survivor) []SurvivorDTO {
	survivors := make([]SurvivorDTO, len(s))
	for i, v := range s {
		survivors[i] = dtoFromDomain(v)
	}
	return survivors
}

type ctxSettlementIdKey string

const SettlementIdKey ctxSettlementIdKey = "settlementId"

func settlementIdExtractor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		settlementIdString := chi.URLParam(r, "id")
		if settlementIdString != "" {
			settlementId, convErr := strconv.Atoi(settlementIdString)
			if convErr != nil {
				web.MakeJsonResponse(w, http.StatusBadRequest, "settlement id should be a number")
				return
			}
			ctx := context.WithValue(r.Context(), SettlementIdKey, settlementId)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
