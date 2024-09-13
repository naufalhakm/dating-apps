package controllers

import (
	"encoding/json"
	"go-dating-test/app/response"
	"go-dating-test/app/services"
	"net/http"
	"strconv"
)

type UserController interface {
	GetUserRecomendation(w http.ResponseWriter, r *http.Request)
}

type UserControllerImpl struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) GetUserRecomendation(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageNum := 1
	limitSize := 5

	if page != "" {
		parsedPage, err := strconv.Atoi(page)
		if err == nil && parsedPage > 0 {
			pageNum = parsedPage
		}
	}

	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err == nil && parsedLimit > 0 {
			limitSize = parsedLimit
		}
	}

	result, pagination, errCuss := controller.UserService.GetUserRecomendation(r.Context(), userID, pageNum, limitSize)

	if errCuss != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errCuss.StatusCode)
		json.NewEncoder(w).Encode(errCuss)
		return
	}

	type Response struct {
		Users      interface{} `json:"users"`
		Pagination interface{} `json:"pagination"`
	}

	var responses Response
	responses.Users = result
	responses.Pagination = pagination

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get data users recomendation.", responses)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}
