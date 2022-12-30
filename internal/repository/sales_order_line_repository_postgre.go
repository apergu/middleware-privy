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

type SalesOrderLineRepositoryPostgre struct {
	pool *pgxpool.Pool
}

func NewSalesOrderLineRepositoryPostgre(pool *pgxpool.Pool) *SalesOrderLineRepositoryPostgre {
	return &SalesOrderLineRepositoryPostgre{
		pool: pool,
	}
}

func (c *SalesOrderLineRepositoryPostgre) query(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) ([]entity.SalesOrderLine, error) {
	rows, err := cmd.Query(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderLineRepositoryPostgre.query",
				"src":   "cmd.Query",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when query",
			"SalesOrderLineRepositoryPostgre.query.Query",
			nil,
		)
	}
	defer rows.Close()

	result := make([]entity.SalesOrderLine, 0)
	for rows.Next() {
		data := entity.SalesOrderLine{}

		err := rows.Scan(
			&data.ID,
			&data.SalesOrderHeaderId,
			&data.ProductID,
			&data.ProductName,
			&data.Quantity,
			&data.RateItem,
			&data.TaxRate,
			&data.Subtotal,
			&data.Grandtotal,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "SalesOrderLineRepositoryPostgre.query",
					"src":   "rows.Scan",
					"query": query,
					"args":  args,
				}).
				Error(err)

			return nil, pgxerror.FromPgxError(
				err,
				"Something went wrong when scan",
				"SalesOrderLineRepositoryPostgre.query.Scan",
			)
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *SalesOrderLineRepositoryPostgre) queryOne(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.SalesOrderLine, error) {
	data := entity.SalesOrderLine{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.SalesOrderHeaderId,
			&data.ProductID,
			&data.ProductName,
			&data.Quantity,
			&data.RateItem,
			&data.TaxRate,
			&data.Subtotal,
			&data.Grandtotal,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderLineRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.SalesOrderLine{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"SalesOrderLineRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *SalesOrderLineRepositoryPostgre) buildFilter(filter SalesOrderLineFilter) (string, []interface{}) {
	condBuilder := &strings.Builder{}
	conds := make([]string, 0, 4) // set for 2 capacity is posible max filter
	condArgs := make([]interface{}, 0, 4)

	if filter.HeaderId > 0 {
		condArgs = append(condArgs, filter.HeaderId)
		idx := "$" + strconv.Itoa(len(condArgs))
		conds = append(conds, "sales_order_lines.sales_order_header_id = "+idx+"")
	}

	if len(conds) > 0 {
		condBuilder.WriteString("where ")
		condBuilder.WriteString(strings.Join(conds, " and "))
	}

	return condBuilder.String(), condArgs
}

func (c *SalesOrderLineRepositoryPostgre) buildSort(sort string) string {
	switch sort {
	case "newest":
		return `order by sales_order_lines.created_at desc`
	}

	return `order by sales_order_lines.updated_at desc`
}

func (c *SalesOrderLineRepositoryPostgre) Find(ctx context.Context, filter SalesOrderLineFilter, limit, skip int64, tx pgx.Tx) ([]entity.SalesOrderLine, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var limits, skips string
	cond, args := c.buildFilter(filter)

	order := c.buildSort(filter.Sort)

	query := `select
		sales_order_lines.id,
		sales_order_lines.sales_order_header_id,
		sales_order_lines.product_id,
		sales_order_lines.product_name,
		sales_order_lines.quantity,
		sales_order_lines.rate_item,
		sales_order_lines.tax_rate,
		sales_order_lines.subtotal,
		sales_order_lines.grandtotal,
		sales_order_lines.created_by,
		sales_order_lines.created_at,
		sales_order_lines.updated_by,
		sales_order_lines.updated_at
	from
		sales_order_lines
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

func (c *SalesOrderLineRepositoryPostgre) Count(ctx context.Context, filter SalesOrderLineFilter, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	cond, args := c.buildFilter(filter)

	query := `select
		count(sales_order_lines.id)
	from
		sales_order_lines
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
			"SalesOrderLineRepositoryPostgre.Count.Scan",
		)
	}

	return data, nil
}

func (c *SalesOrderLineRepositoryPostgre) FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.SalesOrderLine, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
		sales_order_lines.id,
		sales_order_lines.sales_order_header_id,
		sales_order_lines.product_id,
		sales_order_lines.product_name,
		sales_order_lines.quantity,
		sales_order_lines.rate_item,
		sales_order_lines.tax_rate,
		sales_order_lines.subtotal,
		sales_order_lines.grandtotal,
		sales_order_lines.created_by,
		sales_order_lines.created_at,
		sales_order_lines.updated_by,
		sales_order_lines.updated_at
	from
		sales_order_lines
	where
		sales_order_lines.id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *SalesOrderLineRepositoryPostgre) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (c *SalesOrderLineRepositoryPostgre) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (c *SalesOrderLineRepositoryPostgre) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (c *SalesOrderLineRepositoryPostgre) FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.SalesOrderLine, error) {
	if tx == nil {
		return entity.SalesOrderLine{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"SalesOrderLineRepositoryPostgre.FindOneByIdForUpdate",
			nil,
		)
	}
	var cmd sqlcommand.Command = tx

	query := `select
		sales_order_lines.id,
		sales_order_lines.sales_order_header_id,
		sales_order_lines.product_id,
		sales_order_lines.product_name,
		sales_order_lines.quantity,
		sales_order_lines.rate_item,
		sales_order_lines.tax_rate,
		sales_order_lines.subtotal,
		sales_order_lines.grandtotal,
		sales_order_lines.created_by,
		sales_order_lines.created_at,
		sales_order_lines.updated_by,
		sales_order_lines.updated_at
	from
		sales_order_lines
	where
		sales_order_lines.id = $1
	limit 1
	FOR UPDATE`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *SalesOrderLineRepositoryPostgre) Create(ctx context.Context, order entity.SalesOrderLine, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var id int64

	query := `insert into sales_order_lines (
		sales_order_header_id,
		product_id,
		product_name,
		quantity,
		rate_item,
		tax_rate,
		subtotal,
		grandtotal,
		created_by, created_at, updated_by, updated_at
	) values (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		,$11 ,$12
	) RETURNING id`

	err := cmd.
		QueryRow(
			ctx,
			query,
			order.SalesOrderHeaderId,
			order.ProductID,
			order.ProductName,
			order.Quantity,
			order.RateItem,
			order.TaxRate,
			order.Subtotal,
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
				"at":    "SalesOrderLineRepositoryPostgre.create",
				"query": query,
			}).
			Error(err)

		return 0, pgxerror.FromPgxError(err, "", "SalesOrderLineRepositoryPostgre.Create")
	}

	return id, nil
}

func (c *SalesOrderLineRepositoryPostgre) Update(ctx context.Context, id int64, order entity.SalesOrderLine, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update sales_order_lines
	set
		product_id = $1,
		product_name = $2,
		quantity = $3,
		rate_item = $4,
		tax_rate = $5,
		subtotal = $6,
		grandtotal = $7,
		updated_by = $8,
		updated_at = $9
	where
		id = $10`

	_, err := cmd.Exec(
		ctx,
		query,
		order.ProductID,
		order.ProductName,
		order.Quantity,
		order.RateItem,
		order.TaxRate,
		order.Subtotal,
		order.Grandtotal,
		order.UpdatedBy,
		order.UpdatedAt,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "SalesOrderLineRepositoryPostgre.Update")
	}

	return nil
}

func (c *SalesOrderLineRepositoryPostgre) Delete(ctx context.Context, id int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from sales_order_lines where id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "SalesOrderLineRepositoryPostgre.Delete")
	}

	return nil
}

func (c *SalesOrderLineRepositoryPostgre) DeleteByHeader(ctx context.Context, headerId int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from sales_order_lines where sales_order_header_id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		headerId,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "SalesOrderLineRepositoryPostgre.DeleteByHeader")
	}

	return nil
}
