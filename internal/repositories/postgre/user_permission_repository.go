package postgre

import (
	"context"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userPermissionRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewUserPermissionRepository(context context.Context, pool *pgxpool.Pool) interfaces.UserPermissionRepository {
	return &userPermissionRepository{context, pool}
}

func (u userPermissionRepository) Replace(models []models.UserPermission) error {
	tx, err := u.connection.BeginTx(u.context, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(u.context)
		} else {
			tx.Commit(u.context)
		}
	}()

	for _, model := range models {
		_, err = tx.Exec(
			u.context,
			"INSERT INTO ln_user_permissions (user_id, action_id, product_type) "+
				"VALUES ($1, $2, $3) "+
				"ON CONFLICT (id) DO UPDATE "+
				"SET user_id = $1, action_id = $2, product_type = $3;",
			model.UserID, model.ActionID, model.ProductType,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u userPermissionRepository) FindALL(userId int) ([]models.UserPermission, error) {
	var model = models.UserPermission{}
	var models []models.UserPermission
	rows, err := u.connection.Query(
		u.context,
		"SELECT id, user_id, action_id, product_type, created_at, updated_at, deleted_at FROM ln_user_permissions WHERE user_id=$1",
		userId,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&model)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}
