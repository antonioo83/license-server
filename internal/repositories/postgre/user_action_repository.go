package postgre

import (
	"context"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userActionRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewUserActionRepository(context context.Context, pool *pgxpool.Pool) interfaces.UserActionRepository {
	return &userActionRepository{context, pool}
}

func (u userActionRepository) FindALL() (map[string]models.UserAction, error) {
	var model = models.UserAction{}
	models := make(map[string]models.UserAction)
	rows, err := u.connection.Query(
		u.context,
		"SELECT id, action, description, created_at, updated_at, deleted_at FROM ln_user_actions LIMIT 100",
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&model.ID, &model.Action, &model.Description, &model.CreatedAt, &model.UpdatedAt, &model.DeletedAt)
		if err != nil {
			return nil, err
		}
		models[model.Action] = model
	}

	return models, nil
}
