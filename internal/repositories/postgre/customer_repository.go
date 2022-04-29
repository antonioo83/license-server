package postgre

import (
	"context"
	"errors"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type customerRepository struct {
	context    context.Context
	connection *pgxpool.Pool
	licenseRep interfaces.LicenseRepository
}

func NewCustomerRepository(context context.Context, pool *pgxpool.Pool, licenseRep interfaces.LicenseRepository) interfaces.CustomerRepository {
	return &customerRepository{context, pool, licenseRep}
}

func (c customerRepository) Replace(userId int, customer models.Customer, licenses []models.Licence) error {
	tx, err := c.connection.BeginTx(c.context, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() error {
		if err != nil {
			return tx.Rollback(c.context)
		} else {
			return tx.Commit(c.context)
		}
	}()

	var customerId int
	err = c.connection.QueryRow(
		c.context,
		`INSERT INTO 
			ln_customers(user_id, code, type, title, inn, description) 
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (user_id, code) DO UPDATE 
			SET
			  type = excluded.type,
			  title = excluded.title,
			  inn = excluded.inn,
			  description = excluded.description
			RETURNING id`,
		userId, &customer.Code, &customer.Type, &customer.Title, &customer.Inn, &customer.Description,
	).Scan(&customerId)
	if err != nil {
		return err
	}

	err = c.licenseRep.DeleteAll(customerId)
	if err != nil {
		return err
	}

	err = c.licenseRep.MultipleReplace(customerId, licenses)

	return err
}

func (c customerRepository) Delete(userId int, code string) error {
	//TODO implement me
	panic("implement me")
}

func (c customerRepository) FindFull(userId int, customerCode string, licenseCode string) (*models.Customer, error) {
	var err error
	var rows pgx.Rows
	sqlBody :=
		`SELECT 
  		  c.id, c.user_id, c.code, c.type, c.title, c.inn, c.description,
  		  l.product_type, l.callback_url, l.count, l.license_key, l.registration_at, 
  		  l.activation_at, l.expiration_at, l.duration, l.description, l.code
		FROM 
  		  ln_customers c
		LEFT JOIN ln_licenses l ON l.customer_id=c.id`
	if licenseCode != "" {
		rows, err = c.connection.Query(
			c.context,
			sqlBody+` WHERE c.user_id=$1 AND c.code=$2 AND l.code=$3 AND c.deleted_at IS NULL`,
			userId, customerCode, licenseCode,
		)
	} else {
		rows, err = c.connection.Query(
			c.context,
			sqlBody+` WHERE c.user_id=$1 AND c.code=$2 AND c.deleted_at IS NULL`,
			userId, customerCode,
		)
	}
	if err != nil {
		return nil, err
	}

	customers, err := getCustomerModels(rows)
	if err != nil {
		return nil, err
	}

	for _, customer := range customers {
		return &customer, nil
	}

	return &models.Customer{}, nil
}

func getCustomerModels(rows pgx.Rows) (map[int]models.Customer, error) {
	var customers = make(map[int]models.Customer)
	var model models.Customer
	var licence models.Licence
	var customer = models.Customer{}
	lastUserId := 0
	for rows.Next() {
		err := rows.Scan(
			&model.ID, &model.UserID, &model.Code, &model.Type, &model.Title, &model.Inn, &model.Description,
			&licence.ProductType, &licence.CallbackUrl, &licence.Count, &licence.LicenseKey, &licence.RegistrationAt,
			&licence.ActivationAt, &licence.ExpirationAt, &licence.Duration, &licence.Description, &licence.Code,
		)
		if err != nil {
			return nil, err
		}

		if lastUserId != model.ID {
			customer = models.Customer{}
			customer.ID = model.ID
			customer.UserID = model.UserID
			customer.Code = model.Code
			customer.Title = model.Title
			customer.Type = model.Type
			customer.Inn = model.Inn
			customer.Description = model.Description
		}
		lastUserId = model.ID

		customer.Licenses = append(
			customer.Licenses,
			models.Licence{
				ID:             licence.ID,
				Code:           licence.Code,
				CustomerId:     licence.CustomerId,
				ProductType:    licence.ProductType,
				CallbackUrl:    licence.CallbackUrl,
				Count:          licence.Count,
				LicenseKey:     licence.LicenseKey,
				RegistrationAt: licence.RegistrationAt,
				ActivationAt:   licence.ActivationAt,
				ExpirationAt:   licence.ExpirationAt,
				Duration:       licence.Duration,
				Description:    licence.Description,
			},
		)

		customers[customer.ID] = customer
	}

	return customers, nil
}

func (c customerRepository) FindByCode(userId int, code string) (*models.Customer, error) {
	var model models.Customer
	err := c.connection.QueryRow(
		c.context,
		`SELECT 
			   id, user_id, code, type, title, inn, description
			 FROM 
			   ln_customers 
			 WHERE user_id=$1 AND code=$2 AND deleted_at IS NULL`,
		userId, code,
	).Scan(&model.ID, &model.UserID, &model.Code, &model.Type, &model.Title, &model.Inn, &model.Description)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (c customerRepository) IsInDatabase(userId int, code string) (bool, error) {
	model, err := c.FindByCode(userId, code)

	return !(model == nil), err
}
