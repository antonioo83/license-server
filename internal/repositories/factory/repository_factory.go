package factory

import (
	"context"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/antonioo83/license-server/internal/repositories/postgre"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewUserRepository(context context.Context, pool *pgxpool.Pool) interfaces.UserRepository {
	return postgre.NewUserRepository(context, pool)
}

func NewUserPermissionRepository(context context.Context, pool *pgxpool.Pool) interfaces.UserPermissionRepository {
	return postgre.NewUserPermissionRepository(context, pool)
}

func NewUserActionRepository(context context.Context, pool *pgxpool.Pool) interfaces.UserActionRepository {
	return postgre.NewUserActionRepository(context, pool)
}
