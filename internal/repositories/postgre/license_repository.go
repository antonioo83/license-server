package postgre

import (
	"context"
	"errors"
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

func (l licenseRepository) Delete(customerId int, code string) error {
	_, err := l.connection.Exec(
		l.context,
		"DELETE FROM ln_licenses WHERE customer_id=$1 AND code=$2",
		customerId, code,
	)

	return err
}

func (l licenseRepository) FindByCode(code string) (*models.Licence, error) {
	var model models.Licence
	err := l.connection.QueryRow(
		l.context,
		`SELECT 
			   id, customer_id, code, product_type, callback_url, count, license_key, registration_at, activation_at, 
			   expiration_at, duration, description 
			 FROM 
			   ln_licenses 
			 WHERE code=$1`,
		code,
	).Scan(&model.ID, &model.CustomerId, &model.Code, &model.ProductType, &model.CallbackUrl, &model.Count, &model.LicenseKey,
		&model.RegistrationAt, &model.ActivationAt, &model.ExpirationAt, &model.Duration, &model.Description,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (l licenseRepository) IsInDatabase(code string) (bool, error) {
	model, err := l.FindByCode(code)

	return !(model == nil), err
}
