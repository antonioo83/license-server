package factory

import (
	"context"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/antonioo83/license-server/internal/repositories/postgre"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewCustomerRepository(context context.Context, pool *pgxpool.Pool, licenseRep interfaces.LicenseRepository) interfaces.CustomerRepository {
	return postgre.NewCustomerRepository(context, pool, licenseRep)
}

func NewLicenseRepository(context context.Context, pool *pgxpool.Pool) interfaces.LicenseRepository {
	return postgre.NewLicenseRepository(context, pool)
}
