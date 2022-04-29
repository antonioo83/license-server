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
			customerId, model.Code, model.ProductType, model.CallbackURL, model.Count, model.LicenseKey, model.ActivationAt,
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

func (l licenseRepository) UpdateCallbackOptions(licenseId int, isSentCallback uint, callbackAttempts uint) error {
	_, err := l.connection.Query(
		l.context,
		`UPDATE 
			   ln_licenses
			 SET 
			   is_sent_callback=$1,
			   callback_attempts=$2
			 WHERE
			   id=$3`,
		&isSentCallback, &callbackAttempts, &licenseId,
	)
	return err
}

func (l licenseRepository) FindAllExpired(maxAttempts uint, limit uint, offset uint) ([]models.Licence, error) {
	var model = models.Licence{}
	var models []models.Licence
	rows, err := l.connection.Query(
		l.context,
		`SELECT 
			   l.id, l.customer_id, l.code, l.product_type, l.callback_url, l.callback_attempts, l.count, l.license_key, l.registration_at, l.activation_at, 
			   l.expiration_at, l.duration, l.description, c.id, c.code, u.id, u.auth_token 
			 FROM 
			   ln_licenses l
			 LEFT JOIN ln_customers c ON c.id=l.customer_id	
			 LEFT JOIN ln_users u ON u.id=c.user_id
			 WHERE
			   callback_attempts < $1 AND is_sent_callback=0 AND expiration_at < NOW() AND u.deleted_at IS NULL
			 LIMIT $2 OFFSET $3`,
		maxAttempts, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&model.ID, &model.CustomerId, &model.Code, &model.ProductType, &model.CallbackURL, &model.CallbackAttempts,
			&model.Count, &model.LicenseKey, &model.RegistrationAt, &model.ActivationAt, &model.ExpirationAt, &model.Duration,
			&model.Description, &model.Customer.ID, &model.Customer.Code, &model.Customer.User.ID, &model.Customer.User.AuthToken,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
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
	).Scan(&model.ID, &model.CustomerId, &model.Code, &model.ProductType, &model.CallbackURL, &model.Count, &model.LicenseKey,
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
