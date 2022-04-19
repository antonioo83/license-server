package postgre

import (
	"context"
	"errors"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewUserRepository(context context.Context, pool *pgxpool.Pool) interfaces.UserRepository {
	return &userRepository{context, pool}
}

func (u userRepository) Save(model models.User) error {
	_, err := u.connection.Query(
		u.context,
		"INSERT INTO ln_users(code, role, title, auth_token, description)VALUES ($1, $2, $3, $4, $5, NOW())",
		&model.Code, &model.Role, &model.Title, &model.AuthToken, &model.Description,
	)
	return err
}

func (u userRepository) Update(model models.User) error {
	_, err := u.connection.Query(
		u.context,
		"UPDATE ln_users SET role=$2, title=$3, auth_token=$4, description=$5, updated_at=NOW()) WHERE code=$1 AND deleted_at IS NULL",
		&model.Code, &model.Role, &model.Title, &model.AuthToken, &model.Description,
	)
	return err
}

func (u userRepository) Delete(code string) error {
	_, err := u.connection.Query(u.context, "UPDATE ln_users SET deleted_at=NOW() WHERE code=$1 AND deleted_at IS NULL", code)

	return err
}

func (u userRepository) FindByCode(code string) (*models.User, error) {
	var model models.User
	err := u.connection.QueryRow(
		u.context,
		"SELECT id, code, role, title, auth_token, description, created_at, updated_at, deleted_at FROM ln_users WHERE code=$1 AND deleted_at IS NULL",
		code,
	).Scan(&model)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (u userRepository) FindALL(limit int, offset int) (*map[string]models.User, error) {
	var model = models.User{}
	models := make(map[string]models.User)
	rows, err := u.connection.Query(
		u.context,
		"SELECT * FROM ln_users WHERE deleted_at IS NULL LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&model)
		if err != nil {
			return nil, err
		}
		models[model.Code] = model
	}

	return &models, nil
}

func (u userRepository) IsInDatabase(code string) (bool, error) {
	model, err := u.FindByCode(code)

	return !(model == nil), err
}
