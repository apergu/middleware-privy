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

type SalesOrderHeaderRepositoryPostgre struct {
	pool *pgxpool.Pool
}

func NewSalesOrderHeaderRepositoryPostgre(pool *pgxpool.Pool) *SalesOrderHeaderRepositoryPostgre {
	return &SalesOrderHeaderRepositoryPostgre{
		pool: pool,
	}
}

func (c *SalesOrderHeaderRepositoryPostgre) query(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) ([]entity.SalesOrderHeader, error) {
	rows, err := cmd.Query(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderHeaderRepositoryPostgre.query",
				"src":   "cmd.Query",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when query",
			"SalesOrderHeaderRepositoryPostgre.query.Query",
			nil,
		)
	}
	defer rows.Close()

	result := make([]entity.SalesOrderHeader, 0)
	for rows.Next() {
		data := entity.SalesOrderHeader{}

		err := rows.Scan(
			&data.ID,
			&data.OrderNumber,
			&data.CustomerID,
			&data.CustomerName,
			&data.Subtotal,
			&data.Tax,
			&data.Grandtotal,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "SalesOrderHeaderRepositoryPostgre.query",
					"src":   "rows.Scan",
					"query": query,
					"args":  args,
				}).
				Error(err)

			return nil, pgxerror.FromPgxError(
				err,
				"Something went wrong when scan",
				"SalesOrderHeaderRepositoryPostgre.query.Scan",
			)
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *SalesOrderHeaderRepositoryPostgre) queryOne(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.SalesOrderHeader, error) {
	data := entity.SalesOrderHeader{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.OrderNumber,
			&data.CustomerID,
			&data.CustomerName,
			&data.Subtotal,
			&data.Tax,
			&data.Grandtotal,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderHeaderRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.SalesOrderHeader{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"SalesOrderHeaderRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *SalesOrderHeaderRepositoryPostgre) buildFilter(filter SalesOrderHeaderFilter) (string, []interface{}) {
	condBuilder := &strings.Builder{}
	conds := make([]string, 0, 4) // set for 2 capacity is posible max filter
	condArgs := make([]interface{}, 0, 4)

	if len(conds) > 0 {
		condBuilder.WriteString("where ")
		condBuilder.WriteString(strings.Join(conds, " and "))
	}

	return condBuilder.String(), condArgs
}

func (c *SalesOrderHeaderRepositoryPostgre) buildSort(sort string) string {
	switch sort {
	case "newest":
		return `order by sales_order_headers.created_at desc`
	}

	return `order by sales_order_headers.updated_at desc`
}

func (c *SalesOrderHeaderRepositoryPostgre) Find(ctx context.Context, filter SalesOrderHeaderFilter, limit, skip int64, tx pgx.Tx) ([]entity.SalesOrderHeader, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var limits, skips string
	cond, args := c.buildFilter(filter)

	order := c.buildSort(filter.Sort)

	query := `select
		sales_order_headers.id,
		sales_order_headers."order_number",
		sales_order_headers.customer_id,
		sales_order_headers.customer_name,
		sales_order_headers.subtotal,
		sales_order_headers.tax,
		sales_order_headers.grandtotal,
		sales_order_headers.created_by,
		sales_order_headers.created_at,
		sales_order_headers.updated_by,
		sales_order_headers.updated_at
	from
		sales_order_headers
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

func (c *SalesOrderHeaderRepositoryPostgre) Count(ctx context.Context, filter SalesOrderHeaderFilter, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	cond, args := c.buildFilter(filter)

	query := `select
		count(sales_order_headers.id)
	from
		sales_order_headers
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
			"SalesOrderHeaderRepositoryPostgre.Count.Scan",
		)
	}

	return data, nil
}

func (c *SalesOrderHeaderRepositoryPostgre) FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.SalesOrderHeader, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
		sales_order_headers.id,
		sales_order_headers."order_number",
		sales_order_headers.customer_id,
		sales_order_headers.customer_name,
		sales_order_headers.subtotal,
		sales_order_headers.tax,
		sales_order_headers.grandtotal,
		sales_order_headers.created_by,
		sales_order_headers.created_at,
		sales_order_headers.updated_by,
		sales_order_headers.updated_at
	from
		sales_order_headers
	where
		sales_order_headers.id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *SalesOrderHeaderRepositoryPostgre) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (c *SalesOrderHeaderRepositoryPostgre) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (c *SalesOrderHeaderRepositoryPostgre) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (c *SalesOrderHeaderRepositoryPostgre) FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.SalesOrderHeader, error) {
	if tx == nil {
		return entity.SalesOrderHeader{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"SalesOrderHeaderRepositoryPostgre.FindOneByIdForUpdate",
			nil,
		)
	}
	var cmd sqlcommand.Command = tx

	query := `select
		sales_order_headers.id,
		sales_order_headers."order_number",
		sales_order_headers.customer_id,
		sales_order_headers.customer_name,
		sales_order_headers.subtotal,
		sales_order_headers.tax,
		sales_order_headers.grandtotal,
		sales_order_headers.created_by,
		sales_order_headers.created_at,
		sales_order_headers.updated_by,
		sales_order_headers.updated_at
	from
		sales_order_headers
	where
		sales_order_headers.id = $1
	limit 1
	FOR UPDATE`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *SalesOrderHeaderRepositoryPostgre) Create(ctx context.Context, order entity.SalesOrderHeader, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var id int64

	query := `insert into sales_order_headers (
		order_number,
		customer_id,
		customer_name,
		subtotal,
		tax,
		grandtotal,
		created_by, created_at, updated_by, updated_at
	) values (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	) RETURNING id`

	err := cmd.
		QueryRow(
			ctx,
			query,
			order.OrderNumber,
			order.CustomerID,
			order.CustomerName,
			order.Subtotal,
			order.Tax,
			order.Grandtotal,
			order.CreatedBy,
			order.CreatedAt,
			order.UpdatedBy,
			order.UpdatedAt,
		).
		Scan(&id)

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderHeaderRepositoryPostgre.create",
				"query": query,
			}).
			Error(err)

		return 0, pgxerror.FromPgxError(err, "", "SalesOrderHeaderRepositoryPostgre.Create")
	}

	return id, nil
}

func (c *SalesOrderHeaderRepositoryPostgre) Update(ctx context.Context, id int64, order entity.SalesOrderHeader, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update sales_order_headers
	set
		order_number = $1,
		customer_id = $2,
		customer_name = $3,
		subtotal = $4,
		tax = $5,
		grandtotal = $6,
		updated_by = $7,
		updated_at = $8
	where
		id = $9`

	_, err := cmd.Exec(
		ctx,
		query,
		order.OrderNumber,
		order.CustomerID,
		order.CustomerName,
		order.Subtotal,
		order.Tax,
		order.Grandtotal,
		order.UpdatedBy,
		order.UpdatedAt,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "SalesOrderHeaderRepositoryPostgre.Update")
	}

	return nil
}

func (c *SalesOrderHeaderRepositoryPostgre) Delete(ctx context.Context, id int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from sales_order_headers where id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "SalesOrderHeaderRepositoryPostgre.Delete")
	}

	return nil
}
