package postgre

import (
	"context"
	"errors"
	"github.com/antonioo83/license-server/internal/handlers"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepository struct {
	context       context.Context
	connection    *pgxpool.Pool
	permissionRep interfaces.UserPermissionRepository
}

func NewUserRepository(context context.Context, pool *pgxpool.Pool, permissionRep interfaces.UserPermissionRepository) interfaces.UserRepository {
	return &userRepository{context, pool, permissionRep}
}

func (u userRepository) Save(user models.User, permissions []models.UserPermission) error {
	tx, err := u.connection.BeginTx(u.context, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() error {
		if err != nil {
			return tx.Rollback(u.context)
		} else {
			return tx.Commit(u.context)
		}
	}()

	var lastInsertId int
	err = u.connection.QueryRow(
		u.context,
		"INSERT INTO ln_users(code, role, title, auth_token, description)VALUES ($1, $2, $3, $4, $5) RETURNING id",
		&user.Code, &user.Role, &user.Title, &user.AuthToken, &user.Description,
	).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	err = u.permissionRep.MultipleInsert(lastInsertId, permissions)

	return err
}

func (u userRepository) Update(model models.User, permissions []models.UserPermission) error {
	tx, err := u.connection.BeginTx(u.context, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() error {
		if err != nil {
			return tx.Rollback(u.context)
		} else {
			return tx.Commit(u.context)
		}
	}()

	_, err = u.connection.Exec(
		u.context,
		"UPDATE ln_users SET role=$1, title=$2, auth_token=$3, description=$4, updated_at=NOW() WHERE code=$5 AND deleted_at IS NULL",
		&model.Role, &model.Title, &model.AuthToken, &model.Description, &model.Code,
	)
	if err != nil {
		return err
	}

	err = u.permissionRep.Delete(model.ID)
	if err != nil {
		return err
	}

	return u.permissionRep.MultipleInsert(model.ID, permissions)
}

func (u userRepository) Delete(code string) error {
	_, err := u.connection.Exec(u.context, "UPDATE ln_users SET deleted_at=NOW() WHERE code=$1 AND deleted_at IS NULL", code)

	return err
}

func (u userRepository) FindByCode(code string) (*models.User, error) {
	var model models.User
	err := u.connection.QueryRow(
		u.context,
		"SELECT id, code, role, title, auth_token, description, created_at FROM ln_users WHERE code=$1 AND deleted_at IS NULL",
		code,
	).Scan(&model.ID, &model.Code, &model.Role, &model.Title, &model.AuthToken, &model.Description, &model.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (u userRepository) FindByToken(code string) (*models.User, error) {
	var model models.User
	err := u.connection.QueryRow(
		u.context,
		"SELECT id, code, role, title, auth_token, description, created_at FROM ln_users WHERE auth_token=$1 AND deleted_at IS NULL",
		code,
	).Scan(&model.ID, &model.Code, &model.Role, &model.Title, &model.AuthToken, &model.Description, &model.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (u userRepository) FindALL(limit int, offset int) (*[]handlers.UserResponse, error) {
	var users = make(map[int]models.User)
	var user models.User
	var permission models.UserPermission
	var action models.UserAction
	rows, err := u.connection.Query(
		u.context,
		"SELECT "+
			"u.id, u.code, u.role, u.title, u.description, u.created_at, u.updated_at, "+
			"p.user_id, p.productType, a.action "+
			"FROM ln_users u "+
			"LEFT JOIN ln_permissions p ON p.userId=u.id "+
			"LEFT JOIN ln_actions a ON a.permissionId=p.id "+
			"WHERE deleted_at IS NULL "+
			"ORDER BY "+
			"LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	isNew := true
	for rows.Next() {
		err = rows.Scan(
			&user.ID, &user.Code, &user.Title, &user.Description, &user.CreatedAt, &user.UpdatedAt,
			&permission.UserID, &permission.ProductType, &action.Action,
		)
		if err != nil {
			return nil, err
		}

		user.Permissions = append(
			user.Permissions,
			models.UserPermission{
				UserID:      permission.UserID,
				ProductType: permission.ProductType,
				Action:      models.UserAction{Action: action.Action},
			},
		)

	}

	return &models, nil
}

func (u userRepository) IsInDatabase(code string) (bool, error) {
	model, err := u.FindByCode(code)

	return !(model == nil), err
}
