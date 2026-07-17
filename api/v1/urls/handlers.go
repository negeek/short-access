package urls

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"
	urlservice "github.com/negeek/short-access/service/v1/url"
	"github.com/negeek/short-access/utils"
)

// Handler serves the url endpoints, delegating the work to the url service.
type Handler struct {
	urls *urlservice.Service
}

func NewHandler(urls *urlservice.Service) *Handler {
	return &Handler{urls: urls}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	newUrl, err := utils.DecodeBody[urlservice.Url](r)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, ok := v1middlewares.UserID(r.Context())
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest, "Something went wrong. Try again", nil)
		return
	}

	created, err := h.urls.Shorten(r.Context(), userID, &newUrl)
	if err != nil {
		utils.RespondError(w, err)
		return
	}
	utils.JsonResponse(w, true, http.StatusCreated, "Successfully shortened url", created)
}

func (h *Handler) UrlExpiry(w http.ResponseWriter, r *http.Request) {
	detail, err := utils.DecodeBody[DateTimeExpiryDetail](r)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, ok := v1middlewares.UserID(r.Context())
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest, "Something went wrong. Try again", nil)
		return
	}

	updated, err := h.urls.SetExpiry(r.Context(), userID, detail.UrlId, detail.TimeUnit, detail.TimeValue)
	if err != nil {
		utils.RespondError(w, err)
		return
	}
	utils.JsonResponse(w, true, http.StatusOK, "Successfully set url expiry", updated)
}

func (h *Handler) CustomUrl(w http.ResponseWriter, r *http.Request) {
	newUrl, err := utils.DecodeBody[urlservice.Url](r)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, ok := v1middlewares.UserID(r.Context())
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest, "Something went wrong. Try again", nil)
		return
	}

	created, err := h.urls.CreateCustom(r.Context(), userID, &newUrl)
	if err != nil {
		utils.RespondError(w, err)
		return
	}
	utils.JsonResponse(w, true, http.StatusCreated, "Successfully created custom url", created)
}

func (h *Handler) UpdateDeleteUrl(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, ok := v1middlewares.UserID(r.Context())
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest, "Something went wrong. Try again", nil)
		return
	}

	switch r.Method {
	case http.MethodPatch:
		// Load the current record, overlay the incoming fields, then save.
		current, err := h.urls.GetByID(r.Context(), userID, id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
		if err := utils.DecodeBodyInto(r, current); err != nil {
			utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
			return
		}
		updated, err := h.urls.Save(r.Context(), current)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
		utils.JsonResponse(w, true, http.StatusOK, "Successfully updated url", updated)

	case http.MethodPut:
		if _, err := h.urls.GetByID(r.Context(), userID, id); err != nil {
			utils.RespondError(w, err)
			return
		}
		replacement := &urlservice.Url{Id: id}
		if err := utils.DecodeBodyInto(r, replacement); err != nil {
			utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
			return
		}
		updated, err := h.urls.Save(r.Context(), replacement)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
		utils.JsonResponse(w, true, http.StatusOK, "Successfully updated url", updated)

	case http.MethodDelete:
		if err := h.urls.Delete(r.Context(), userID, id); err != nil {
			utils.RespondError(w, err)
			return
		}
		utils.JsonResponse(w, true, http.StatusNoContent, "Successfully deleted url", nil)
	}
}

func (h *Handler) UrlFilter(w http.ResponseWriter, r *http.Request) {
	userID, ok := v1middlewares.UserID(r.Context())
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest, "Something went wrong. Try again", nil)
		return
	}

	results, err := h.urls.List(r.Context(), userID, r.URL.Query())
	if err != nil {
		utils.RespondError(w, err)
		return
	}
	utils.JsonResponse(w, true, http.StatusOK, "", results)
}

func (h *Handler) UrlRedirect(w http.ResponseWriter, r *http.Request) {
	target, err := h.urls.Redirect(r.Context(), mux.Vars(r)["slug"])
	if err != nil {
		utils.RespondError(w, err)
		return
	}
	http.Redirect(w, r, target.OriginalUrl, http.StatusTemporaryRedirect)
}
