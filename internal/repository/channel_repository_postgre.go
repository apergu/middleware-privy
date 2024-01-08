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

type ChannelRepositoryPostgre struct {
	pool *pgxpool.Pool
}

func NewChannelRepositoryPostgre(pool *pgxpool.Pool) *ChannelRepositoryPostgre {
	return &ChannelRepositoryPostgre{
		pool: pool,
	}
}

func (c *ChannelRepositoryPostgre) query(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) ([]entity.Channel, error) {
	rows, err := cmd.Query(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelRepositoryPostgre.query",
				"src":   "cmd.Query",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when query",
			"ChannelRepositoryPostgre.query.Query",
			nil,
		)
	}
	defer rows.Close()

	result := make([]entity.Channel, 0)
	for rows.Next() {
		data := entity.Channel{}

		err := rows.Scan(
			&data.ID,
			&data.MerchantID,
			&data.ChannelCode,
			&data.ChannelID,
			&data.ChannelName,
			&data.Address,
			&data.Email,
			&data.PhoneNo,
			&data.State,
			&data.City,
			&data.ZipCode,
			&data.CustomerInternalID,
			&data.MerchantID,
			&data.ChannelInternalID,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "ChannelRepositoryPostgre.query",
					"src":   "rows.Scan",
					"query": query,
					"args":  args,
				}).
				Error(err)

			return nil, pgxerror.FromPgxError(
				err,
				"Something went wrong when scan",
				"ChannelRepositoryPostgre.query.Scan",
			)
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *ChannelRepositoryPostgre) queryOne(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.Channel, error) {
	data := entity.Channel{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.MerchantID,
			&data.ChannelCode,
			&data.ChannelID,
			&data.ChannelName,
			&data.Address,
			&data.Email,
			&data.PhoneNo,
			&data.State,
			&data.City,
			&data.ZipCode,
			&data.CustomerInternalID,
			&data.MerchantID,
			&data.ChannelInternalID,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.Channel{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"ChannelRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *ChannelRepositoryPostgre) buildFilter(filter ChannelFilter) (string, []interface{}) {
	condBuilder := &strings.Builder{}
	conds := make([]string, 0, 4) // set for 2 capacity is posible max filter
	condArgs := make([]interface{}, 0, 4)

	if filter.ChannelID != nil {
		condArgs = append(condArgs, *filter.ChannelID)
		idx := "$" + strconv.Itoa(len(condArgs))

		conds = append(conds, "channel_id = "+idx)
	}

	if len(conds) > 0 {
		condBuilder.WriteString("where ")
		condBuilder.WriteString(strings.Join(conds, " and "))
	}

	return condBuilder.String(), condArgs
}

func (c *ChannelRepositoryPostgre) buildSort(sort string) string {
	switch sort {
	case "newest":
		return `order by channels.created_at desc`
	}

	return `order by channels.updated_at desc`
}

func (c *ChannelRepositoryPostgre) Find(ctx context.Context, filter ChannelFilter, limit, skip int64, tx pgx.Tx) ([]entity.Channel, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var limits, skips string
	cond, args := c.buildFilter(filter)

	order := c.buildSort(filter.Sort)

	query := `select
		channels.id,
		channels.merchant_id,
		channels.channel_code,
		channels.channel_id,
		channels.channel_name,
		channels."address",
		channels.email,
		channels.phone_no,
		channels.state,
		channels.city,
		channels.zip_code,
		channels.customer_internalid,
		channels.merchant_internalid,
		channels.channel_internalid,
		channels.created_by,
		channels.created_at,
		channels.updated_by,
		channels.updated_at
	from
		channels
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

func (c *ChannelRepositoryPostgre) Count(ctx context.Context, filter ChannelFilter, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	cond, args := c.buildFilter(filter)

	query := `select
		count(channels.id)
	from
		channels
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
			"ChannelRepositoryPostgre.Count.Scan",
		)
	}

	return data, nil
}

func (c *ChannelRepositoryPostgre) FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Channel, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
		channels.id,
		channels.merchant_id,
		channels.channel_code,
		channels.channel_id,
		channels.channel_name,
		channels."address",
		channels."email",
		channels.phone_no,
		channels.state,
		channels."city",
		channels.zip_code,
		channels.customer_internalid,
		channels.merchant_internalid,
		channels.channel_internalid,
		channels.created_by,
		channels.created_at,
		channels.updated_by,
		channels.updated_at
	from
		channels
	where
		channels.id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *ChannelRepositoryPostgre) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (c *ChannelRepositoryPostgre) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (c *ChannelRepositoryPostgre) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (c *ChannelRepositoryPostgre) FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Channel, error) {
	if tx == nil {
		return entity.Channel{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"ChannelRepositoryPostgre.FindOneByIdForUpdate",
			nil,
		)
	}
	var cmd sqlcommand.Command = tx

	query := `select
		channels."id",
		channels.merchant_id,
		channels.channel_code,
		channels.channel_id,
		channels.channel_name,
		channels."address",
		channels."email",
		channels.phone_no,
		channels.state,
		channels."city",
		channels.zip_code,
		channels.customer_internalid,
		channels.merchant_internalid,
		channels.channel_internalid,
		channels.created_by,
		channels.created_at,
		channels.updated_by,
		channels.updated_at
	from
		channels
	where
		channels.id = $1
	limit 1
	FOR UPDATE`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *ChannelRepositoryPostgre) Create(ctx context.Context, channel entity.Channel, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	// Check for duplicate channel_id
	duplicateQuery := `SELECT id FROM channels WHERE channel_id = $1 LIMIT 1`
	var existingID int64
	err := cmd.
		QueryRow(ctx, duplicateQuery, channel.ChannelID).
		Scan(&existingID)

	if err == nil {
		// Duplicate entry found
		return 0, fmt.Errorf("duplicate entry with channel_id %s", channel.ChannelID)
	} else if err != pgx.ErrNoRows {
		// An error occurred while checking for duplicates
		return 0, pgxerror.FromPgxError(err, "", "ChannelRepositoryPostgre.Create")
	}

	var id int64
	query := `insert into channels (
		merchant_id,
		channel_code,
		channel_id,
		channel_name,
		"address",
		email,
		phone_no,
		"state",
		"city",
		"zip_code",
		customer_internalid,
		merchant_internalid,
		channel_internalid,
		created_by, created_at, updated_by, updated_at
	) values (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		,$11, $12, $13, $14, $15, $16, $17
	) RETURNING id`

	err = cmd.
		QueryRow(
			ctx,
			query,
			channel.MerchantID,
			channel.ChannelCode,
			channel.ChannelID,
			channel.ChannelName,
			channel.Address,
			channel.Email,
			channel.PhoneNo,
			channel.State,
			channel.City,
			channel.ZipCode,
			channel.CustomerInternalID,
			channel.MerchantInternalID,
			channel.ChannelInternalID,
			channel.CreatedBy,
			channel.CreatedAt,
			channel.UpdatedBy,
			channel.UpdatedAt,
		).
		Scan(&id)

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelRepositoryPostgre.create",
				"query": query,
			}).
			Error(err)

		return 0, pgxerror.FromPgxError(err, "", "ChannelRepositoryPostgre.Create")
	}

	return id, nil
}

func (c *ChannelRepositoryPostgre) Update(ctx context.Context, id int64, channel entity.Channel, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update channels
	set
		channel_id = $1,
		channel_name = $2,
		address = $3,
		email = $4,
		phone_no = $5,
		state = $6,
		"city" = $7,
		"zip_code" = $8,
		channel_code = $12,
		customer_internalid = $13,
		merchant_internalid = $14,
		channel_internalid = $15,
		updated_by = $9,
		updated_at = $10
	where
		id = $11`

	_, err := cmd.Exec(
		ctx,
		query,
		channel.ChannelID,
		channel.ChannelName,
		channel.Address,
		channel.Email,
		channel.PhoneNo,
		channel.State,
		channel.City,
		channel.ZipCode,
		channel.UpdatedBy,
		channel.UpdatedAt,
		id,
		channel.ChannelCode,
		channel.CustomerInternalID,
		channel.MerchantInternalID,
		channel.ChannelInternalID,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "ChannelRepositoryPostgre.Update")
	}

	return nil
}

func (c *ChannelRepositoryPostgre) Delete(ctx context.Context, id int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from channels where id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "ChannelRepositoryPostgre.Delete")
	}

	return nil
}
