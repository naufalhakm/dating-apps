package services

import (
	"context"
	"database/sql"
	"go-dating-test/app/params"
	"go-dating-test/app/repositories"
	"go-dating-test/app/response"
)

type UserService interface {
	GetUserRecomendation(ctx context.Context, id string, page, pageSize int) ([]*params.UserResponse, *params.Pagination, *response.CustomError)
}

type UserServiceImpl struct {
	DB             *sql.DB
	UserRepository repositories.UserRepository
}

func NewUserService(db *sql.DB, userRepository repositories.UserRepository) UserService {
	return &UserServiceImpl{
		DB:             db,
		UserRepository: userRepository,
	}
}

func (service *UserServiceImpl) GetUserRecomendation(ctx context.Context, id string, page, pageSize int) ([]*params.UserResponse, *params.Pagination, *response.CustomError) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to database errors: %s", err.Error())
	}
	defer func() {
		err := recover()
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	userLogin, err := service.UserRepository.GetUserDetail(ctx, tx, id)
	if err != nil {
		return nil, nil, response.BadRequestErrorWithAdditionalInfo("users not found.")
	}

	var pagination = new(params.Pagination)
	pagination.Page = (page - 1) * pageSize
	pagination.PageSize = pageSize

	users, err := service.UserRepository.GetAllUserRecomendation(ctx, tx, userLogin, pagination)
	if err != nil {
		return nil, nil, response.BadRequestError()
	}
	pagination.Page = page
	pagination.PageCount = (pagination.TotalCount + pagination.PageSize - 1) / pagination.PageSize

	var result []*params.UserResponse
	for _, user := range users {

		var preferencesResponse = new(params.PreferencesResponse)
		preferencesResponse.PreferredAgeRange = user.Preferences.PreferredAgeRange
		preferencesResponse.PreferredGender = user.Preferences.PreferredGender
		preferencesResponse.MaxDistanceKm = user.Preferences.MaxDistanceKm

		result = append(result, &params.UserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Age:         user.Age,
			Gender:      user.Gender,
			Latitude:    user.Latitude,
			Longitude:   user.Longitude,
			Interests:   user.Interests,
			Preferences: preferencesResponse,
			LastActive:  user.LastActive,
		})
	}

	return result, pagination, nil
}
