package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/pkg/pgxerror"
	"gitlab.com/mohamadikbal/project-privy/pkg/sqlcommand"
	"gitlab.com/rteja-library3/rapperror"
)

type CustomerRepositoryPostgre struct {
	pool *pgxpool.Pool
}

func NewCustomerRepositoryPostgre(pool *pgxpool.Pool) *CustomerRepositoryPostgre {
	return &CustomerRepositoryPostgre{
		pool: pool,
	}
}

func (c *CustomerRepositoryPostgre) query(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) ([]entity.Customer, error) {
	rows, err := cmd.Query(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerRepositoryPostgre.query",
				"src":   "cmd.Query",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when query",
			"CustomerRepositoryPostgre.query.Query",
			nil,
		)
	}
	defer rows.Close()

	result := make([]entity.Customer, 0)
	for rows.Next() {
		data := entity.Customer{}

		err := rows.Scan(
			&data.ID,
			&data.CustomerID,
			&data.CustomerType,
			&data.CustomerName,
			&data.FirstName,
			&data.LastName,
			&data.Email,
			&data.PhoneNo,
			&data.Address,
			&data.CRMLeadID,
			&data.EnterprisePrivyID,
			&data.CustomerInternalID,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "CustomerRepositoryPostgre.query",
					"src":   "rows.Scan",
					"query": query,
					"args":  args,
				}).
				Error(err)

			return nil, pgxerror.FromPgxError(
				err,
				"Something went wrong when scan",
				"CustomerRepositoryPostgre.query.Scan",
			)
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *CustomerRepositoryPostgre) queryOne(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.Customer, error) {
	data := entity.Customer{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.CustomerID,
			&data.CustomerType,
			&data.CustomerName,
			&data.FirstName,
			&data.LastName,
			&data.Email,
			&data.PhoneNo,
			&data.Address,
			&data.CRMLeadID,
			&data.EnterprisePrivyID,
			&data.CustomerInternalID,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.Customer{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"CustomerRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *CustomerRepositoryPostgre) buildFilter(filter CustomerFilter) (string, []interface{}) {
	condBuilder := &strings.Builder{}
	conds := make([]string, 0, 4) // set for 2 capacity is posible max filter
	condArgs := make([]interface{}, 0, 4)

	if filter.EnterprisePrivyID != nil {
		condArgs = append(condArgs, *filter.EnterprisePrivyID)
		idx := "$" + strconv.Itoa(len(condArgs))

		conds = append(conds, "enterprise_privy_id = "+idx)
	}

	if len(conds) > 0 {
		condBuilder.WriteString("where ")
		condBuilder.WriteString(strings.Join(conds, " and "))
	}

	return condBuilder.String(), condArgs
}

func (c *CustomerRepositoryPostgre) buildSort(sort string) string {
	switch sort {
	case "newest":
		return `order by customers.created_at desc`
	}

	return `order by customers.updated_at desc`
}

func (c *CustomerRepositoryPostgre) Find(ctx context.Context, filter CustomerFilter, limit, skip int64, tx pgx.Tx) ([]entity.Customer, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var limits, skips string
	cond, args := c.buildFilter(filter)

	order := c.buildSort(filter.Sort)

	query := `select
		customers.id,
		customers.customer_id,
		customers.customer_type,
		customers.customer_name,
		customers.first_name,
		customers.last_name,
		customers.email,
		customers.phone_no,
		customers."address",
		customers."crm_lead_id",
		customers."enterprise_privy_id",
		customers."customer_internalid",
		customers.created_by,
		customers.created_at,
		customers.updated_by,
		customers.updated_at
	from
		customers
	%s
	%s
	%s
	%s`

	if limit > 0 {
		args = append(args, limit)
		limits = "limit $" + strconv.Itoa(len(args))
	}

	if skip > 0 {
		args = append(args, skip)
		skips = "offset $" + strconv.Itoa(len(args))
	}

	return c.query(ctx, cmd, fmt.Sprintf(query, cond, order, limits, skips), args...)
}

func (c *CustomerRepositoryPostgre) Count(ctx context.Context, filter CustomerFilter, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	cond, args := c.buildFilter(filter)

	query := `select
		count(customers.id)
	from
		customers
	%s`

	var data int64
	err := cmd.
		QueryRow(ctx, fmt.Sprintf(query, cond), args...).
		Scan(
			&data,
		)
	if err != nil {
		return 0, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"CustomerRepositoryPostgre.Count.Scan",
		)
	}

	return data, nil
}

func (c *CustomerRepositoryPostgre) FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Customer, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
		customers.id,
		customers.customer_id,
		customers.customer_type,
		customers.customer_name,
		customers.first_name,
		customers.last_name,
		customers.email,
		customers.phone_no,
		customers."address",
		customers."crm_lead_id",
		customers."enterprise_privy_id",
		customers."customer_internalid",
		customers.created_by,
		customers.created_at,
		customers.updated_by,
		customers.updated_at
	from
		customers
	where
		customers.id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *CustomerRepositoryPostgre) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (c *CustomerRepositoryPostgre) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (c *CustomerRepositoryPostgre) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (c *CustomerRepositoryPostgre) FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Customer, error) {
	if tx == nil {
		return entity.Customer{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"CustomerRepositoryPostgre.FindOneByIdForUpdate",
			nil,
		)
	}
	var cmd sqlcommand.Command = tx

	query := `select
		customers.id,
		customers.customer_id,
		customers.customer_type,
		customers.customer_name,
		customers.first_name,
		customers.last_name,
		customers.email,
		customers.phone_no,
		customers."address",
		customers."crm_lead_id",
		customers."enterprise_privy_id",
		customers."customer_internalid",
		customers.created_by,
		customers.created_at,
		customers.updated_by,
		customers.updated_at
	from
		customers
	where
		customers.id = $1
	limit 1
	FOR UPDATE`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *CustomerRepositoryPostgre) Create(ctx context.Context, cust entity.Customer, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var id int64

	query := `insert into customers (
		customer_id,
		customer_type,
		customer_name,
		first_name,
		last_name,
		email,
		phone_no,
		"address",
		"crm_lead_id",
		"enterprise_privy_id",
		"address_1",
		"npwp",
		"state",
		"city",
		"zip_code",
		"customer_internalid",
		created_by, created_at, updated_by, updated_at
	) values (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		,$11, $12 ,$13, $14, $15, $16, $17, $18, $19, $20
	) RETURNING id`

	err := cmd.
		QueryRow(
			ctx,
			query,
			cust.CustomerID,
			cust.CustomerType,
			cust.CustomerName,
			cust.FirstName,
			cust.LastName,
			cust.Email,
			cust.PhoneNo,
			cust.Address,
			cust.CRMLeadID,
			cust.EnterprisePrivyID,
			cust.Address1,
			cust.NPWP,
			cust.State,
			cust.City,
			cust.ZipCode,
			cust.CustomerInternalID,
			cust.CreatedBy,
			cust.CreatedAt,
			cust.UpdatedBy,
			cust.UpdatedAt,
		).
		Scan(&id)

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerRepositoryPostgre.create",
				"query": query,
			}).
			Error(err)

		return 0, pgxerror.FromPgxError(err, "", "CustomerRepositoryPostgre.Create")
	}

	return id, nil
}

func (c *CustomerRepositoryPostgre) Update(ctx context.Context, id int64, cust entity.Customer, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update customers
	set
		customer_type = $1,
		customer_name = $2,
		first_name = $3,
		last_name = $4,
		email = $5,
		phone_no = $6,
		"address" = $7,
		"crm_lead_id" = $8,
		"enterprise_privy_id" = $9,
		"address_1" = $13,
		"npwp" = $14,
		"state" = $15,
		"city" = $16,
		"zip_code" = $17,
		"customer_internalid" = $18,
		updated_by = $10,
		updated_at = $11
	where
		id = $12`

	_, err := cmd.Exec(
		ctx,
		query,
		// cust.CustomerID,
		cust.CustomerType,
		cust.CustomerName,
		cust.FirstName,
		cust.LastName,
		cust.Email,
		cust.PhoneNo,
		cust.Address,
		cust.CRMLeadID,
		cust.EnterprisePrivyID,
		cust.UpdatedBy,
		cust.UpdatedAt,
		id,
		cust.Address1,
		cust.NPWP,
		cust.State,
		cust.City,
		cust.ZipCode,
		cust.CustomerInternalID,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "CustomerRepositoryPostgre.Update")
	}

	return nil
}

func (c *CustomerRepositoryPostgre) Delete(ctx context.Context, id int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from customers where id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "CustomerRepositoryPostgre.Delete")
	}

	return nil
}
