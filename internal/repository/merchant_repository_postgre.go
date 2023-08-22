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

type MerchantRepositoryPostgre struct {
	pool *pgxpool.Pool
}

func NewMerchantRepositoryPostgre(pool *pgxpool.Pool) *MerchantRepositoryPostgre {
	return &MerchantRepositoryPostgre{
		pool: pool,
	}
}

func (c *MerchantRepositoryPostgre) query(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) ([]entity.Merchant, error) {
	rows, err := cmd.Query(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantRepositoryPostgre.query",
				"src":   "cmd.Query",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when query",
			"MerchantRepositoryPostgre.query.Query",
			nil,
		)
	}
	defer rows.Close()

	result := make([]entity.Merchant, 0)
	for rows.Next() {
		data := entity.Merchant{}

		err := rows.Scan(
			&data.ID,
			&data.CustomerID,
			&data.EnterpriseID,
			&data.MerchantID,
			&data.MerchantName,
			&data.Address,
			&data.Email,
			&data.PhoneNo,
			&data.State,
			&data.City,
			&data.ZipCode,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "MerchantRepositoryPostgre.query",
					"src":   "rows.Scan",
					"query": query,
					"args":  args,
				}).
				Error(err)

			return nil, pgxerror.FromPgxError(
				err,
				"Something went wrong when scan",
				"MerchantRepositoryPostgre.query.Scan",
			)
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *MerchantRepositoryPostgre) queryOne(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.Merchant, error) {
	data := entity.Merchant{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.CustomerID,
			&data.EnterpriseID,
			&data.MerchantID,
			&data.MerchantName,
			&data.Address,
			&data.Email,
			&data.PhoneNo,
			&data.State,
			&data.City,
			&data.ZipCode,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.Merchant{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"MerchantRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *MerchantRepositoryPostgre) buildFilter(filter MerchantFilter) (string, []interface{}) {
	condBuilder := &strings.Builder{}
	conds := make([]string, 0, 4) // set for 2 capacity is posible max filter
	condArgs := make([]interface{}, 0, 4)

	if len(conds) > 0 {
		condBuilder.WriteString("where ")
		condBuilder.WriteString(strings.Join(conds, " and "))
	}

	return condBuilder.String(), condArgs
}

func (c *MerchantRepositoryPostgre) buildSort(sort string) string {
	switch sort {
	case "newest":
		return `order by merchants.created_at desc`
	}

	return `order by merchants.updated_at desc`
}

func (c *MerchantRepositoryPostgre) Find(ctx context.Context, filter MerchantFilter, limit, skip int64, tx pgx.Tx) ([]entity.Merchant, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var limits, skips string
	cond, args := c.buildFilter(filter)

	order := c.buildSort(filter.Sort)

	query := `select
		merchants.id,
		merchants.customer_id,
		merchants.enterprise_id,
		merchants.merchant_id,
		merchants.merchant_name,
		merchants."address",
		merchants.email,
		merchants.phone_no,
		merchants.state,
		merchants.city,
		merchants.zip_code,
		merchants.created_by,
		merchants.created_at,
		merchants.updated_by,
		merchants.updated_at
	from
		merchants
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

func (c *MerchantRepositoryPostgre) Count(ctx context.Context, filter MerchantFilter, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	cond, args := c.buildFilter(filter)

	query := `select
		count(merchants.id)
	from
		merchants
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
			"MerchantRepositoryPostgre.Count.Scan",
		)
	}

	return data, nil
}

func (c *MerchantRepositoryPostgre) FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Merchant, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
		merchants.id,
		merchants.customer_id,
		merchants.enterprise_id,
		merchants.merchant_id,
		merchants.merchant_name,
		merchants."address",
		merchants.email,
		merchants.phone_no,
		merchants.state,
		merchants.city,
		merchants.zip_code,
		merchants.created_by,
		merchants.created_at,
		merchants.updated_by,
		merchants.updated_at
	from
		merchants
	where
		merchants.id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *MerchantRepositoryPostgre) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (c *MerchantRepositoryPostgre) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (c *MerchantRepositoryPostgre) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (c *MerchantRepositoryPostgre) FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Merchant, error) {
	if tx == nil {
		return entity.Merchant{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"MerchantRepositoryPostgre.FindOneByIdForUpdate",
			nil,
		)
	}
	var cmd sqlcommand.Command = tx

	query := `select
		merchants.id,
		merchants.customer_id,
		merchants.enterprise_id,
		merchants.merchant_id,
		merchants.merchant_name,
		merchants."address",
		merchants.email,
		merchants.phone_no,
		merchants.state,
		merchants.city,
		merchants.zip_code,
		merchants.created_by,
		merchants.created_at,
		merchants.updated_by,
		merchants.updated_at
	from
		merchants
	where
		merchants.id = $1
	limit 1
	FOR UPDATE`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *MerchantRepositoryPostgre) Create(ctx context.Context, merchant entity.Merchant, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var id int64
	query := `insert into merchants (
		customer_id,
		enterprise_id,
		merchant_id,
		merchant_name,
		"address",
		email,
		phone_no,
		"state",
		"city",
		"zip_code",
		created_by, created_at, updated_by, updated_at
	) values (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		,$11, $12 ,$13, $14
	) RETURNING id`

	err := cmd.
		QueryRow(
			ctx,
			query,
			merchant.CustomerID,
			merchant.EnterpriseID,
			merchant.MerchantID,
			merchant.MerchantName,
			merchant.Address,
			merchant.Email,
			merchant.PhoneNo,
			merchant.State,
			merchant.City,
			merchant.ZipCode,
			merchant.CreatedBy,
			merchant.CreatedAt,
			merchant.UpdatedBy,
			merchant.UpdatedAt,
		).
		Scan(&id)

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantRepositoryPostgre.create",
				"query": query,
			}).
			Error(err)

		return 0, pgxerror.FromPgxError(err, "", "MerchantRepositoryPostgre.Create")
	}

	return id, nil
}

func (c *MerchantRepositoryPostgre) Update(ctx context.Context, id int64, merchant entity.Merchant, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update merchants
	set
		merchant_id = $1,
		merchant_name = $2,
		address = $3,
		email = $4,
		phone_no = $5,
		state = $6,
		"city" = $7,
		"zip_code" = $8,
		updated_by = $9,
		updated_at = $10
	where
		id = $11`

	_, err := cmd.Exec(
		ctx,
		query,
		merchant.MerchantID,
		merchant.MerchantName,
		merchant.Address,
		merchant.Email,
		merchant.PhoneNo,
		merchant.State,
		merchant.City,
		merchant.ZipCode,
		merchant.UpdatedBy,
		merchant.UpdatedAt,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "MerchantRepositoryPostgre.Update")
	}

	return nil
}

func (c *MerchantRepositoryPostgre) Delete(ctx context.Context, id int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from merchants where id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "MerchantRepositoryPostgre.Delete")
	}

	return nil
}
