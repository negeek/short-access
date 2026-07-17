package apikeys

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"
	apikeyservice "github.com/negeek/short-access/service/v1/apikey"
	"github.com/negeek/short-access/utils"
)

// Handler serves the API-key management endpoints. These sit behind JWT auth, so
// a user signs in to manage the keys their application will use.
type Handler struct {
	keys *apikeyservice.Service
}

func NewHandler(keys *apikeyservice.Service) *Handler {
	return &Handler{keys: keys}
}

// createRequest is the body of a create-key request. expire_at is optional; when
// omitted the key never expires.
type createRequest struct {
	Name     string     `json:"name"`
	ExpireAt *time.Time `json:"expire_at"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := v1middlewares.UserID(r.Context())
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest, "Something went wrong. Try again", nil)
		return
	}

	req, err := utils.DecodeBody[createRequest](r)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	raw, record, err := h.keys.Create(r.Context(), userID, req.Name, req.ExpireAt)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// The raw key is only ever returned here, so tell the caller to save it now.
	utils.JsonResponse(w, true, http.StatusCreated, "API key created. Copy it now, it will not be shown again.", map[string]interface{}{
		"api_key": raw,
		"key":     record,
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := v1middlewares.UserID(r.Context())
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest, "Something went wrong. Try again", nil)
		return
	}

	keys, err := h.keys.List(r.Context(), userID)
	if err != nil {
		utils.RespondError(w, err)
		return
	}
	utils.JsonResponse(w, true, http.StatusOK, "", keys)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := v1middlewares.UserID(r.Context())
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest, "Something went wrong. Try again", nil)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.keys.Delete(r.Context(), userID, id); err != nil {
		utils.RespondError(w, err)
		return
	}
	utils.JsonResponse(w, true, http.StatusOK, "API key deleted", nil)
}

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	userID, ok := v1middlewares.UserID(r.Context())
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest, "Something went wrong. Try again", nil)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.keys.Revoke(r.Context(), userID, id); err != nil {
		utils.RespondError(w, err)
		return
	}
	utils.JsonResponse(w, true, http.StatusOK, "API key revoked", nil)
}
