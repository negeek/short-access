package users

import (
	"net/http"

	usersvc "github.com/negeek/short-access/service/v1/user"
	"github.com/negeek/short-access/utils"
)

// Handler serves the user endpoints, delegating the work to the user service.
type Handler struct {
	users *usersvc.Service
}

func NewHandler(users *usersvc.Service) *Handler {
	return &Handler{users: users}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	newUser, err := utils.DecodeBody[usersvc.User](r)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := h.users.SignUp(r.Context(), &newUser)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.JsonResponse(w, true, http.StatusCreated, "Successfully Joined", map[string]interface{}{
		"email":        newUser.Email,
		"access_token": token,
	})
}

func (h *Handler) NewToken(w http.ResponseWriter, r *http.Request) {
	oldUser, err := utils.DecodeBody[usersvc.User](r)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := h.users.NewToken(r.Context(), &oldUser)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.JsonResponse(w, true, http.StatusCreated, "Token created Successfully", map[string]interface{}{
		"email":        oldUser.Email,
		"access_token": token,
	})
}
