package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"

	"middleware/internal/entity"
	"middleware/pkg/pgxerror"
	"middleware/pkg/sqlcommand"
)

type ApplicationRepositoryPostgre struct {
	pool *pgxpool.Pool
}

func NewApplicationRepositoryPostgre(pool *pgxpool.Pool) *ApplicationRepositoryPostgre {
	return &ApplicationRepositoryPostgre{
		pool: pool,
	}
}

func (c *ApplicationRepositoryPostgre) query(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) ([]entity.Application, error) {
	rows, err := cmd.Query(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationRepositoryPostgre.query",
				"src":   "cmd.Query",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when query",
			"ApplicationRepositoryPostgre.query.Query",
			nil,
		)
	}
	defer rows.Close()

	result := make([]entity.Application, 0)
	for rows.Next() {
		data := entity.Application{}

		err := rows.Scan(
			&data.ID,
			&data.ApplicationCode,
			&data.ApplicationID,
			&data.ApplicationName,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "ApplicationRepositoryPostgre.query",
					"src":   "rows.Scan",
					"query": query,
					"args":  args,
				}).
				Error(err)

			return nil, pgxerror.FromPgxError(
				err,
				"Something went wrong when scan",
				"ApplicationRepositoryPostgre.query.Scan",
			)
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *ApplicationRepositoryPostgre) queryOne(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.Application, error) {
	data := entity.Application{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.ApplicationCode,
			&data.ApplicationID,
			&data.ApplicationName,
		)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.Application{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"ApplicationRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *ApplicationRepositoryPostgre) queryOneFind(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.ApplicationFind, error) {
	data := entity.ApplicationFind{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.ApplicationID,
			&data.ApplicationName,
		)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.ApplicationFind{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"ApplicationRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *ApplicationRepositoryPostgre) buildFilter(filter ApplicationFilter) (string, []interface{}) {
	condBuilder := &strings.Builder{}
	conds := make([]string, 0, 4) // set for 2 capacity is posible max filter
	condArgs := make([]interface{}, 0, 4)

	if filter.ApplicationID != nil {
		condArgs = append(condArgs, *filter.ApplicationID)
		idx := "$" + strconv.Itoa(len(condArgs))

		conds = append(conds, "application_id = "+idx)
	}

	if len(conds) > 0 {
		condBuilder.WriteString("where ")
		condBuilder.WriteString(strings.Join(conds, " and "))
	}

	return condBuilder.String(), condArgs
}

func (c *ApplicationRepositoryPostgre) buildSort(sort string) string {
	switch sort {
	case "newest":
		return `order by applications.created_at desc`
	}

	return `order by applications.updated_at desc`
}

func (c *ApplicationRepositoryPostgre) Find(ctx context.Context, filter ApplicationFilter, limit, skip int64, tx pgx.Tx) ([]entity.Application, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var limits, skips string
	cond, args := c.buildFilter(filter)

	order := c.buildSort(filter.Sort)

	query := `select
		applications.id,
		applications.merchant_id,
		applications.application_code,
		applications.application_id,
		applications.application_name,
		applications."address",
		applications.email,
		applications.phone_no,
		applications.state,
		applications.city,
		applications.zip_code,
		applications.customer_internalid,
		applications.merchant_internalid,
		applications.application_internalid,
		applications.created_by,
		applications.created_at,
		applications.updated_by,
		applications.updated_at
	from
		applications
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

func (c *ApplicationRepositoryPostgre) FindByApplicationID(ctx context.Context, enterprisePrivyID string, tx pgx.Tx) (entity.ApplicationFind, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
	applications.id,
	applications.application_id,
	applications.application_name
	from
		applications
	where
		applications.application_id = $1
	limit 1`

	return c.queryOneFind(ctx, cmd, query, enterprisePrivyID)
}

func (c *ApplicationRepositoryPostgre) FindByName(ctx context.Context, enterprisePrivyID string, tx pgx.Tx) (entity.ApplicationFind, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
	applications.id,
	applications.application_id,
	applications.application_name
	from
		applications
	where
		applications.application_name = $1
	limit 1`

	return c.queryOneFind(ctx, cmd, query, enterprisePrivyID)
}

func (c *ApplicationRepositoryPostgre) FindByMerchantID(ctx context.Context, merchantId string, tx pgx.Tx) (entity.Application, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
	applications.id,
		applications.merchant_id,
		applications.application_code,
		applications.application_id,
		applications.application_name,
		applications."address",
		applications."email",
		applications.phone_no,
		applications.state,
		applications."city",
		applications.zip_code,
		applications.customer_internalid,
		applications.merchant_internalid,
		applications.application_internalid,
		applications.created_by,
		applications.created_at,
		applications.updated_by,
		applications.updated_at
	from
		applications
	where
		applications.merchant_id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, merchantId)
}

func (c *ApplicationRepositoryPostgre) Count(ctx context.Context, filter ApplicationFilter, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	cond, args := c.buildFilter(filter)

	query := `select
		count(applications.id)
	from
		applications
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
			"ApplicationRepositoryPostgre.Count.Scan",
		)
	}

	return data, nil
}

func (c *ApplicationRepositoryPostgre) FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Application, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
		applications.id,
		applications.merchant_id,
		applications.application_code,
		applications.application_id,
		applications.application_name,
		applications."address",
		applications."email",
		applications.phone_no,
		applications.state,
		applications."city",
		applications.zip_code,
		applications.customer_internalid,
		applications.merchant_internalid,
		applications.application_internalid,
		applications.created_by,
		applications.created_at,
		applications.updated_by,
		applications.updated_at
	from
		applications
	where
		applications.id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *ApplicationRepositoryPostgre) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (c *ApplicationRepositoryPostgre) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (c *ApplicationRepositoryPostgre) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (c *ApplicationRepositoryPostgre) FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Application, error) {
	if tx == nil {
		return entity.Application{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"ApplicationRepositoryPostgre.FindOneByIdForUpdate",
			nil,
		)
	}
	var cmd sqlcommand.Command = tx

	query := `select
		applications."id",
		applications.merchant_id,
		applications.application_code,
		applications.application_id,
		applications.application_name,
		applications."address",
		applications."email",
		applications.phone_no,
		applications.state,
		applications."city",
		applications.zip_code,
		applications.customer_internalid,
		applications.merchant_internalid,
		applications.application_internalid,
		applications.created_by,
		applications.created_at,
		applications.updated_by,
		applications.updated_at
	from
		applications
	where
		applications.id = $1
	limit 1
	FOR UPDATE`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *ApplicationRepositoryPostgre) Create(ctx context.Context, application entity.Application, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	// 1 merchan now can have multi application code above is not necessary
	var id int64
	query := `insert into applications (
		merchant_id,
		application_code,
		application_id,
		application_name,
		"address",
		email,
		phone_no,
		"state",
		"city",
		"zip_code",
		customer_internalid,
		merchant_internalid,
		application_internalid,
		created_by, created_at, updated_by, updated_at
	) values (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		,$11, $12, $13, $14, $15, $16, $17
	) RETURNING id`

	err := cmd.
		QueryRow(
			ctx,
			query,
			application.EnterpriseID,
			application.ApplicationCode,
			application.ApplicationID,
			application.ApplicationName,
		).
		Scan(&id)

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationRepositoryPostgre.create",
				"query": query,
			}).
			Error(err)

		return 0, pgxerror.FromPgxError(err, "", "ApplicationRepositoryPostgre.Create")
	}

	return id, nil
}

func (c *ApplicationRepositoryPostgre) Update(ctx context.Context, id int64, application entity.Application, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update applications
	set
		application_id = $1,
		application_name = $2,
		address = $3,
		email = $4,
		phone_no = $5,
		state = $6,
		"city" = $7,
		"zip_code" = $8,
		application_code = $12,
		customer_internalid = $13,
		merchant_internalid = $14,
		application_internalid = $15,
		updated_by = $9,
		updated_at = $10
	where
		id = $11`

	_, err := cmd.Exec(
		ctx,
		query,
		application.EnterpriseID,
		application.ApplicationID,
		application.ApplicationName,
		id,
		application.ApplicationCode,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "ApplicationRepositoryPostgre.Update")
	}

	return nil
}

func (c *ApplicationRepositoryPostgre) Delete(ctx context.Context, id int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from applications where id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "ApplicationRepositoryPostgre.Delete")
	}

	return nil
}
