package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"middleware/internal/entity"
	"middleware/pkg/pgxerror"
	"middleware/pkg/sqlcommand"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type CustomerUsageRepositoryPostgre struct {
	pool *pgxpool.Pool
}

func NewCustomerUsageRepositoryPostgre(pool *pgxpool.Pool) *CustomerUsageRepositoryPostgre {
	return &CustomerUsageRepositoryPostgre{
		pool: pool,
	}
}

func (c *CustomerUsageRepositoryPostgre) query(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) ([]entity.CustomerUsage, error) {
	rows, err := cmd.Query(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageRepositoryPostgre.query",
				"src":   "cmd.Query",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when query",
			"CustomerUsageRepositoryPostgre.query.Query",
			nil,
		)
	}
	defer rows.Close()

	result := make([]entity.CustomerUsage, 0)
	for rows.Next() {
		data := entity.CustomerUsage{}

		err := rows.Scan(
			&data.ID,
			&data.CustomerID,
			&data.CustomerName,
			&data.ProductID,
			&data.ProductName,
			&data.TransactionAt,
			&data.Balance,
			&data.BalanceAmount,
			&data.Usage,
			&data.UsageAmount,
			&data.EnterpriseID,
			&data.EnterpriseName,
			&data.ChannelName,
			&data.TrxId,
			&data.ServiceID,
			&data.UnitPrice,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "CustomerUsageRepositoryPostgre.query",
					"src":   "rows.Scan",
					"query": query,
					"args":  args,
				}).
				Error(err)

			return nil, pgxerror.FromPgxError(
				err,
				"Something went wrong when scan",
				"CustomerUsageRepositoryPostgre.query.Scan",
			)
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *CustomerUsageRepositoryPostgre) queryOne(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.CustomerUsage, error) {
	data := entity.CustomerUsage{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.CustomerID,
			&data.CustomerName,
			&data.ProductID,
			&data.ProductName,
			&data.TransactionAt,
			&data.Balance,
			&data.BalanceAmount,
			&data.Usage,
			&data.UsageAmount,
			&data.EnterpriseID,
			&data.EnterpriseName,
			&data.ChannelName,
			&data.TrxId,
			&data.ServiceID,
			&data.UnitPrice,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.CustomerUsage{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"CustomerUsageRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *CustomerUsageRepositoryPostgre) buildFilter(filter CustomerUsageFilter) (string, []interface{}) {
	condBuilder := &strings.Builder{}
	conds := make([]string, 0, 4) // set for 2 capacity is posible max filter
	condArgs := make([]interface{}, 0, 4)

	if len(conds) > 0 {
		condBuilder.WriteString("where ")
		condBuilder.WriteString(strings.Join(conds, " and "))
	}

	return condBuilder.String(), condArgs
}

func (c *CustomerUsageRepositoryPostgre) buildSort(sort string) string {
	switch sort {
	case "newest":
		return `order by customer_usages.created_at desc`
	}

	return `order by customer_usages.updated_at desc`
}

func (c *CustomerUsageRepositoryPostgre) Find(ctx context.Context, filter CustomerUsageFilter, limit, skip int64, tx pgx.Tx) ([]entity.CustomerUsage, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var limits, skips string
	cond, args := c.buildFilter(filter)

	order := c.buildSort(filter.Sort)

	query := `select
		customer_usages.id,
		customer_usages.customer_id,
		customer_usages.customer_name,
		customer_usages.product_id,
		customer_usages.product_name,
		customer_usages.transaction_at,
		customer_usages.balance,
		customer_usages.balance_amount,
		customer_usages.usage,
		customer_usages.usage_amount,
		customer_usages.enterprise_id,
		customer_usages.enterprise_name,
		customer_usages.channel_name,
		customer_usages.trx_id,
		customer_usages.service_id,
		customer_usages.unit_price,
		customer_usages.created_by,
		customer_usages.created_at,
		customer_usages.updated_by,
		customer_usages.updated_at
	from
		customer_usages
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

func (c *CustomerUsageRepositoryPostgre) Count(ctx context.Context, filter CustomerUsageFilter, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	cond, args := c.buildFilter(filter)

	query := `select
		count(customer_usages.id)
	from
		customer_usages
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
			"CustomerUsageRepositoryPostgre.Count.Scan",
		)
	}

	return data, nil
}

func (c *CustomerUsageRepositoryPostgre) FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.CustomerUsage, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
		customer_usages.id,
		customer_usages.customer_id,
		customer_usages.customer_name,
		customer_usages.product_id,
		customer_usages.product_name,
		customer_usages.transaction_at,
		customer_usages.balance,
		customer_usages.balance_amount,
		customer_usages.usage,
		customer_usages.usage_amount,
		customer_usages.enterprise_id,
		customer_usages.enterprise_name,
		customer_usages.channel_name,
		customer_usages.trx_id,
		customer_usages.service_id,
		customer_usages.unit_price,
		customer_usages.created_by,
		customer_usages.created_at,
		customer_usages.updated_by,
		customer_usages.updated_at
	from
		customer_usages
	where
		customer_usages.id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *CustomerUsageRepositoryPostgre) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (c *CustomerUsageRepositoryPostgre) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (c *CustomerUsageRepositoryPostgre) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (c *CustomerUsageRepositoryPostgre) FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.CustomerUsage, error) {
	if tx == nil {
		return entity.CustomerUsage{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"CustomerUsageRepositoryPostgre.FindOneByIdForUpdate",
			nil,
		)
	}
	var cmd sqlcommand.Command = tx

	query := `select
		customer_usages.id,
		customer_usages.customer_id,
		customer_usages.customer_name,
		customer_usages.product_id,
		customer_usages.product_name,
		customer_usages.transaction_at,
		customer_usages.balance,
		customer_usages.balance_amount,
		customer_usages.usage,
		customer_usages.usage_amount,
		customer_usages.enterprise_id,
		customer_usages.enterprise_name,
		customer_usages.channel_name,
		customer_usages.trx_id,
		customer_usages.service_id,
		customer_usages.unit_price,
		customer_usages.created_by,
		customer_usages.created_at,
		customer_usages.updated_by,
		customer_usages.updated_at
	from
		customer_usages
	where
		customer_usages.id = $1
	limit 1
	FOR UPDATE`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *CustomerUsageRepositoryPostgre) Create(ctx context.Context, cust entity.CustomerUsage, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var id int64

	query := `insert into customer_usages (
		customer_id,
		customer_name,
		product_id,
		product_name,
		transaction_at,
		balance,
		balance_amount,
		usage,
		usage_amount,
		enterprise_id,
		enterprise_name,
		channel_name,
		trx_id,
		service_id,
		unit_price,
		created_by, created_at, updated_by, updated_at
	) values (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15, $16, $17, $18, $19
	) RETURNING id`

	err := cmd.
		QueryRow(
			ctx,
			query,
			cust.CustomerID,
			cust.CustomerName,
			cust.ProductID,
			cust.ProductName,
			cust.TransactionAt,
			cust.Balance,
			cust.BalanceAmount,
			cust.Usage,
			cust.UsageAmount,
			cust.EnterpriseID,
			cust.EnterpriseName,
			cust.ChannelName,
			cust.TrxId,
			cust.ServiceID,
			cust.UnitPrice,
			cust.CreatedBy,
			cust.CreatedAt,
			cust.UpdatedBy,
			cust.UpdatedAt,
		).
		Scan(&id)

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageRepositoryPostgre.create",
				"query": query,
			}).
			Error(err)

		return 0, pgxerror.FromPgxError(err, "", "CustomerUsageRepositoryPostgre.Create")
	}

	return id, nil
}

func (c *CustomerUsageRepositoryPostgre) Update(ctx context.Context, id int64, cust entity.CustomerUsage, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update customer_usages
	set
		enterprise_id = $13,
		enterprise_name = $14,
		channel_name = $15,
		trx_id = $16,
		service_id = $17,
		unit_price = $18,
		customer_id = $1,
		customer_name = $2,
		product_id = $3,
		product_name = $4,
		transaction_at = $5,
		balance = $6,
		balance_amount = $7,
		usage = $8,
		usage_amount = $9,
		updated_by = $10,
		updated_at = $11
	where
		id = $12`

	_, err := cmd.Exec(
		ctx,
		query,
		cust.CustomerID,
		cust.CustomerName,
		cust.ProductID,
		cust.ProductName,
		cust.TransactionAt,
		cust.Balance,
		cust.BalanceAmount,
		cust.Usage,
		cust.UsageAmount,
		cust.UpdatedBy,
		cust.UpdatedAt,
		id,
		cust.EnterpriseID,
		cust.EnterpriseName,
		cust.ChannelName,
		cust.TrxId,
		cust.ServiceID,
		cust.UnitPrice,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "CustomerUsageRepositoryPostgre.Update")
	}

	return nil
}

func (c *CustomerUsageRepositoryPostgre) Delete(ctx context.Context, id int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from customer_usages where id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "CustomerUsageRepositoryPostgre.Delete")
	}

	return nil
}
