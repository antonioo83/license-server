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

func (u userRepository) FindFullUser(code string) (*models.User, error) {
	rows, err := u.connection.Query(
		u.context,
		`SELECT 
				u.id, u.code, u.role, u.title, u.description, 
				p.user_id, p.product_type, a.action 
			FROM ln_users u 
			LEFT JOIN ln_user_permissions p ON user_id=u.id 
			LEFT JOIN ln_user_actions a ON a.id=p.action_id
			WHERE 
  				u.code=$1 AND u.deleted_at IS NULL`,
		code,
	)
	if err != nil {
		return nil, err
	}

	users, err := getModels(rows)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		return &user, nil
	}

	return &models.User{}, nil
}

func (u userRepository) FindFullUsers(limit int, offset int) (*map[int]models.User, error) {
	rows, err := u.connection.Query(
		u.context,
		`SELECT 
				u.id, u.code, u.role, u.title, u.description, 
				p.user_id, p.product_type, a.action 
			FROM ln_users u 
			LEFT JOIN ln_user_permissions p ON user_id=u.id 
			LEFT JOIN ln_user_actions a ON a.id=p.action_id
			WHERE 
  				u.deleted_at IS NULL 
			ORDER BY u.id ASC
			LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}

	users, err := getModels(rows)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func getModels(rows pgx.Rows) (map[int]models.User, error) {
	var users = make(map[int]models.User)
	var model models.User
	var permission models.UserPermission
	var action models.UserAction
	var user = models.User{}
	lastUserId := 0
	for rows.Next() {
		err := rows.Scan(
			&model.ID, &model.Code, &model.Role, &model.Title, &model.Description, //&model.CreatedAt, &model.UpdatedAt,
			&permission.UserID, &permission.ProductType, &action.Action,
		)
		if err != nil {
			return nil, err
		}

		if lastUserId != model.ID {
			user = models.User{}
			user.ID = model.ID
			user.Role = model.Role
			user.Code = model.Code
			user.Title = model.Title
			user.Description = model.Description
		}
		lastUserId = model.ID

		user.Permissions = append(
			user.Permissions,
			models.UserPermission{
				UserID:      permission.UserID,
				ProductType: permission.ProductType,
				Action:      models.UserAction{Action: action.Action},
			},
		)

		users[user.ID] = user
	}

	return users, nil
}

func (u userRepository) IsInDatabase(code string) (bool, error) {
	model, err := u.FindByCode(code)

	return !(model == nil), err
}
