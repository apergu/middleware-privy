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

type DivissionRepositoryPostgre struct {
	pool *pgxpool.Pool
}

func NewDivissionRepositoryPostgre(pool *pgxpool.Pool) *DivissionRepositoryPostgre {
	return &DivissionRepositoryPostgre{
		pool: pool,
	}
}

func (c *DivissionRepositoryPostgre) query(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) ([]entity.Divission, error) {
	rows, err := cmd.Query(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "DivissionRepositoryPostgre.query",
				"src":   "cmd.Query",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when query",
			"DivissionRepositoryPostgre.query.Query",
			nil,
		)
	}
	defer rows.Close()

	result := make([]entity.Divission, 0)
	for rows.Next() {
		data := entity.Divission{}

		err := rows.Scan(
			&data.ID,
			&data.ChannelID,
			&data.DivissionID,
			&data.DivissionName,
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
					"at":    "DivissionRepositoryPostgre.query",
					"src":   "rows.Scan",
					"query": query,
					"args":  args,
				}).
				Error(err)

			return nil, pgxerror.FromPgxError(
				err,
				"Something went wrong when scan",
				"DivissionRepositoryPostgre.query.Scan",
			)
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *DivissionRepositoryPostgre) queryOne(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.Divission, error) {
	data := entity.Divission{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.ChannelID,
			&data.DivissionID,
			&data.DivissionName,
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
				"at":    "DivissionRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.Divission{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"DivissionRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *DivissionRepositoryPostgre) buildFilter(filter DivissionFilter) (string, []interface{}) {
	condBuilder := &strings.Builder{}
	conds := make([]string, 0, 4) // set for 2 capacity is posible max filter
	condArgs := make([]interface{}, 0, 4)

	if len(conds) > 0 {
		condBuilder.WriteString("where ")
		condBuilder.WriteString(strings.Join(conds, " and "))
	}

	return condBuilder.String(), condArgs
}

func (c *DivissionRepositoryPostgre) buildSort(sort string) string {
	switch sort {
	case "newest":
		return `order by divissions.created_at desc`
	}

	return `order by divissions.updated_at desc`
}

func (c *DivissionRepositoryPostgre) Find(ctx context.Context, filter DivissionFilter, limit, skip int64, tx pgx.Tx) ([]entity.Divission, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var limits, skips string
	cond, args := c.buildFilter(filter)

	order := c.buildSort(filter.Sort)

	query := `select
		divissions.id,
		divissions.channel_id,
		divissions.divission_id,
		divissions.divission_name,
		divissions."address",
		divissions.email,
		divissions.phone_no,
		divissions.state,
		divissions.city,
		divissions.zip_code,
		divissions.created_by,
		divissions.created_at,
		divissions.updated_by,
		divissions.updated_at
	from
		divissions
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

func (c *DivissionRepositoryPostgre) Count(ctx context.Context, filter DivissionFilter, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	cond, args := c.buildFilter(filter)

	query := `select
		count(divissions.id)
	from
		divissions
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
			"DivissionRepositoryPostgre.Count.Scan",
		)
	}

	return data, nil
}

func (c *DivissionRepositoryPostgre) FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Divission, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
		divissions.id,
		divissions.channel_id,
		divissions.divission_id,
		divissions.divission_name,
		divissions."address",
		divissions."email",
		divissions.phone_no,
		divissions.state,
		divissions."city",
		divissions.zip_code,
		divissions.created_by,
		divissions.created_at,
		divissions.updated_by,
		divissions.updated_at
	from
		divissions
	where
		divissions.id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *DivissionRepositoryPostgre) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (c *DivissionRepositoryPostgre) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (c *DivissionRepositoryPostgre) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (c *DivissionRepositoryPostgre) FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Divission, error) {
	if tx == nil {
		return entity.Divission{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"DivissionRepositoryPostgre.FindOneByIdForUpdate",
			nil,
		)
	}
	var cmd sqlcommand.Command = tx

	query := `select
		divissions."id",
		divissions.channel_id,
		divissions.divission_id,
		divissions.divission_name,
		divissions."address",
		divissions."email",
		divissions.phone_no,
		divissions.state,
		divissions."city",
		divissions.zip_code,
		divissions.created_by,
		divissions.created_at,
		divissions.updated_by,
		divissions.updated_at
	from
		divissions
	where
		divissions.id = $1
	limit 1
	FOR UPDATE`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *DivissionRepositoryPostgre) Create(ctx context.Context, divission entity.Divission, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var id int64
	query := `insert into divissions (
		channel_id,
		divission_id,
		divission_name,
		"address",
		email,
		phone_no,
		"state",
		"city",
		"zip_code",
		created_by, created_at, updated_by, updated_at
	) values (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		,$11, $12 ,$13
	) RETURNING id`

	err := cmd.
		QueryRow(
			ctx,
			query,
			divission.ChannelID,
			divission.DivissionID,
			divission.DivissionName,
			divission.Address,
			divission.Email,
			divission.PhoneNo,
			divission.State,
			divission.City,
			divission.ZipCode,
			divission.CreatedBy,
			divission.CreatedAt,
			divission.UpdatedBy,
			divission.UpdatedAt,
		).
		Scan(&id)

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "DivissionRepositoryPostgre.create",
				"query": query,
			}).
			Error(err)

		return 0, pgxerror.FromPgxError(err, "", "DivissionRepositoryPostgre.Create")
	}

	return id, nil
}

func (c *DivissionRepositoryPostgre) Update(ctx context.Context, id int64, divission entity.Divission, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update divissions
	set
		divission_id = $1,
		divission_name = $2,
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
		divission.DivissionID,
		divission.DivissionName,
		divission.Address,
		divission.Email,
		divission.PhoneNo,
		divission.State,
		divission.City,
		divission.ZipCode,
		divission.UpdatedBy,
		divission.UpdatedAt,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "DivissionRepositoryPostgre.Update")
	}

	return nil
}

func (c *DivissionRepositoryPostgre) Delete(ctx context.Context, id int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from divissions where id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "DivissionRepositoryPostgre.Delete")
	}

	return nil
}
