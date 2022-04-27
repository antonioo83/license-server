package postgre

import (
	"context"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type licenseRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewLicenseRepository(context context.Context, pool *pgxpool.Pool) interfaces.LicenseRepository {
	return &licenseRepository{context, pool}
}

func (l licenseRepository) MultipleReplace(customerId int, models []models.Licence) error {
	b := &pgx.Batch{}
	for _, model := range models {
		b.Queue(
			`INSERT INTO 
					  ln_licenses(customer_id, code, product_type, callback_url, count, license_key, activation_at, expiration_at, duration, description) 
				   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
				   ON CONFLICT (customer_id, code) DO UPDATE 
 				   SET 
					  product_type = excluded.product_type, 
					  callback_url = excluded.callback_url, 
					  count = excluded.count, 
					  license_key = excluded.license_key, 
					  activation_at = excluded.activation_at, 
					  expiration_at = excluded.expiration_at,  
					  duration = excluded.duration, 
					  description = excluded.description`,
			customerId, model.Code, model.ProductType, model.CallbackUrl, model.Count, model.LicenseKey, model.ActivationAt,
			model.ExpirationAt, model.Duration, model.Description,
		)
	}
	r := l.connection.SendBatch(l.context, b)
	_, err := r.Exec()
	if err != nil {
		return err
	}

	return r.Close()
}

func (l licenseRepository) DeleteAll(customerId int) error {
	_, err := l.connection.Exec(
		l.context,
		"DELETE FROM ln_licenses WHERE customer_id=$1",
		customerId,
	)

	return err
}
