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

func (u userPermissionRepository) MultipleInsert(userId int, models []models.UserPermission) error {
	b := &pgx.Batch{}
	for _, model := range models {
		b.Queue(
			"INSERT INTO ln_user_permissions (user_id, action_id, product_type) VALUES ($1, $2, $3)",
			userId, model.ActionID, model.ProductType,
		)
	}
	r := u.connection.SendBatch(u.context, b)
	_, err := r.Exec()
	if err != nil {
		return err
	}

	return r.Close()
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

func (u userPermissionRepository) Delete(userId int) error {
	_, err := u.connection.Exec(u.context, "DELETE FROM ln_user_permissions WHERE user_id=$1", userId)

	return err
}
