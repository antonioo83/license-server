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
