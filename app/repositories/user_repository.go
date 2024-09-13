package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go-dating-test/app/models"
	"go-dating-test/app/params"

	"github.com/lib/pq"
)

type UserRepository interface {
	GetUserDetail(ctx context.Context, tx *sql.Tx, id string) (*models.User, error)
	GetAllUserRecomendation(ctx context.Context, tx *sql.Tx, userLogin *models.User, pagination *params.Pagination) ([]*models.User, error)
}

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) GetUserDetail(ctx context.Context, tx *sql.Tx, id string) (*models.User, error) {
	query := `
		SELECT id, name, age, gender, 
		       ST_X(location::geometry) AS longitude, 
		       ST_Y(location::geometry) AS latitude, 
		       interests, 
		       preferences 
		FROM users 
		WHERE id = $1`

	row := tx.QueryRowContext(ctx, query, id)

	var user models.User
	var interests []string
	var preferencesData []byte

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Age,
		&user.Gender,
		&user.Longitude,
		&user.Latitude,
		pq.Array(&interests),
		&preferencesData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user is not found")
		}
		return nil, err
	}

	var preferences models.Preferences
	if err := json.Unmarshal(preferencesData, &preferences); err != nil {
		return nil, err
	}

	user.Interests = interests
	user.Preferences = &preferences

	return &user, nil
}

func (repository *UserRepositoryImpl) GetAllUserRecomendation(ctx context.Context, tx *sql.Tx, userLogin *models.User, pagination *params.Pagination) ([]*models.User, error) {

	query :=
		`SELECT id, name, age, gender, 
           ST_X(location::geometry) AS longitude, 
           ST_Y(location::geometry) AS latitude, 
           interests, 
           preferences,
           last_active,
           ARRAY_LENGTH(ARRAY(SELECT UNNEST(interests) INTERSECT SELECT UNNEST(ARRAY[$7::text[]])),1) AS matching_interest_count,
		   count(*) over() as total
    FROM users 
    WHERE ST_DWithin(location, ST_MakePoint($1, $2)::geography, $3 * 1000) 
    AND (age BETWEEN $4 AND $5) 
    AND gender = $6
	AND id != $8
    ORDER BY matching_interest_count ASC, last_active DESC
	LIMIT $9 OFFSET $10;`

	var users []*models.User

	rows, err := tx.QueryContext(
		ctx,
		query,
		userLogin.Longitude,
		userLogin.Latitude,
		userLogin.Preferences.MaxDistanceKm,
		userLogin.Preferences.PreferredAgeRange[0],
		userLogin.Preferences.PreferredAgeRange[1],
		userLogin.Preferences.PreferredGender,
		pq.Array(userLogin.Interests),
		userLogin.ID,
		pagination.PageSize,
		pagination.Page,
	)

	if err != nil {
		fmt.Printf("Query error: %v\n", err)
		return nil, err
	}

	for rows.Next() {
		var user models.User
		var interests []string
		var preferencesData []byte
		var matchingInterestCount sql.NullInt64

		err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Gender, &user.Longitude, &user.Latitude, pq.Array(&interests), &preferencesData, &user.LastActive, &matchingInterestCount, &pagination.TotalCount)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			return nil, err
		}

		var preferences models.Preferences
		if err := json.Unmarshal(preferencesData, &preferences); err != nil {
			return nil, err
		}

		user.Interests = interests
		user.Preferences = &preferences

		users = append(users, &user)
	}

	return users, nil
}
