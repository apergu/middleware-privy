package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"middleware/internal/entity"
	"middleware/pkg/pgxerror"
	"middleware/pkg/sqlcommand"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type TransferBalanceRepositoryPostgre struct {
	pool *pgxpool.Pool
}

func NewTransferBalanceRepositoryPostgre(pool *pgxpool.Pool) *TransferBalanceRepositoryPostgre {
	return &TransferBalanceRepositoryPostgre{
		pool: pool,
	}
}

func (c *TransferBalanceRepositoryPostgre) query(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) ([]entity.TransferBalance, error) {
	rows, err := cmd.Query(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TransferBalanceRepositoryPostgre.query",
				"src":   "cmd.Query",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when query",
			"TransferBalanceRepositoryPostgre.query.Query",
			nil,
		)
	}
	defer rows.Close()

	result := make([]entity.TransferBalance, 0)
	for rows.Next() {
		data := entity.TransferBalance{}

		err := rows.Scan(
			&data.ID,
			&data.CustomerId,
			&data.MerchantTo,
			&data.ChannelTo,
			&data.StartDate,
			&data.EndDate,
			&data.Quantity,
			&data.TransferDate,
			&data.IsTrxCreated,
			&data.TrxIdFrom,
			&data.TrxIdTo,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "TransferBalanceRepositoryPostgre.query",
					"src":   "rows.Scan",
					"query": query,
					"args":  args,
				}).
				Error(err)

			return nil, pgxerror.FromPgxError(
				err,
				"Something went wrong when scan",
				"TransferBalanceRepositoryPostgre.query.Scan",
			)
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *TransferBalanceRepositoryPostgre) queryOne(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.TransferBalance, error) {
	data := entity.TransferBalance{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.CustomerId,
			&data.MerchantTo,
			&data.ChannelTo,
			&data.StartDate,
			&data.EndDate,
			&data.Quantity,
			&data.TransferDate,
			&data.IsTrxCreated,
			&data.TrxIdFrom,
			&data.TrxIdTo,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TransferBalanceRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.TransferBalance{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"TransferBalanceRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *TransferBalanceRepositoryPostgre) buildFilter(filter TransferBalanceFilter) (string, []interface{}) {
	condBuilder := &strings.Builder{}
	conds := make([]string, 0, 4) // set for 2 capacity is posible max filter
	condArgs := make([]interface{}, 0, 4)

	if len(conds) > 0 {
		condBuilder.WriteString("where ")
		condBuilder.WriteString(strings.Join(conds, " and "))
	}

	return condBuilder.String(), condArgs
}

func (c *TransferBalanceRepositoryPostgre) buildSort(sort string) string {
	switch sort {
	case "newest":
		return `order by transfer_balances.created_at desc`
	}

	return `order by transfer_balances.updated_at desc`
}

func (c *TransferBalanceRepositoryPostgre) Find(ctx context.Context, filter TransferBalanceFilter, limit, skip int64, tx pgx.Tx) ([]entity.TransferBalance, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var limits, skips string
	cond, args := c.buildFilter(filter)

	order := c.buildSort(filter.Sort)

	query := `select
		transfer_balances."id",
		transfer_balances.merchant_id,
		transfer_balances.transaction_id,
		transfer_balances.enterprise_id,
		transfer_balances."enterprise_name",
		transfer_balances.original_service_id,
		transfer_balances."service_id",
		transfer_balances.service_name,
		transfer_balances.quantity,
		transfer_balances.transaction_date,
		transfer_balances.merchant_code,
		transfer_balances.channel_id,
		transfer_balances.channel_code,
		transfer_balances.customer_internalid,
		transfer_balances.merchant_internalid,
		transfer_balances.channel_internalid,
		transfer_balances.transaction_type,
		transfer_balances.topup_id,
		transfer_balances.created_by,
		transfer_balances.created_at,
		transfer_balances.updated_by,
		transfer_balances.updated_at
	from
		transfer_balances
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

func (c *TransferBalanceRepositoryPostgre) Count(ctx context.Context, filter TransferBalanceFilter, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	cond, args := c.buildFilter(filter)

	query := `select
		count(transfer_balances.id)
	from
		transfer_balances
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
			"TransferBalanceRepositoryPostgre.Count.Scan",
		)
	}

	return data, nil
}

func (c *TransferBalanceRepositoryPostgre) FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.TransferBalance, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
		transfer_balances."id",
		transfer_balances.merchant_id,
		transfer_balances.transaction_id,
		transfer_balances.enterprise_id,
		transfer_balances."enterprise_name",
		transfer_balances.original_service_id,
		transfer_balances."service_id",
		transfer_balances.service_name,
		transfer_balances.quantity,
		transfer_balances.transaction_date,
		transfer_balances.merchant_code,
		transfer_balances.channel_id,
		transfer_balances.channel_code,
		transfer_balances.customer_internalid,
		transfer_balances.merchant_internalid,
		transfer_balances.channel_internalid,
		transfer_balances.transaction_type,
		transfer_balances.topup_id,
		transfer_balances.created_by,
		transfer_balances.created_at,
		transfer_balances.updated_by,
		transfer_balances.updated_at
	from
		transfer_balances
	where
		transfer_balances.id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *TransferBalanceRepositoryPostgre) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (c *TransferBalanceRepositoryPostgre) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (c *TransferBalanceRepositoryPostgre) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (c *TransferBalanceRepositoryPostgre) FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.TransferBalance, error) {
	if tx == nil {
		return entity.TransferBalance{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"TransferBalanceRepositoryPostgre.FindOneByIdForUpdate",
			nil,
		)
	}
	var cmd sqlcommand.Command = tx

	query := `select
		transfer_balances."id",
		transfer_balances.merchant_id,
		transfer_balances.transaction_id,
		transfer_balances.enterprise_id,
		transfer_balances."enterprise_name",
		transfer_balances.original_service_id,
		transfer_balances."service_id",
		transfer_balances.service_name,
		transfer_balances.quantity,
		transfer_balances.transaction_date,
		transfer_balances.merchant_code,
		transfer_balances.channel_id,
		transfer_balances.channel_code,
		transfer_balances.customer_internalid,
		transfer_balances.merchant_internalid,
		transfer_balances.channel_internalid,
		transfer_balances.transaction_type,
		transfer_balances.topup_id,
		transfer_balances.created_by,
		transfer_balances.created_at,
		transfer_balances.updated_by,
		transfer_balances.updated_at
	from
		transfer_balances
	where
		transfer_balances.id = $1
	limit 1
	FOR UPDATE`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *TransferBalanceRepositoryPostgre) Create(ctx context.Context, topup entity.TransferBalance, tx pgx.Tx) (any, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	// topUpUUid := uuid.New().String()

	var id int64
	query := `insert into transfer_balance (
		customer_id,
		transfer_date,
		trx_id_from,
        trx_id_to,
        merchant_to,
		channel_to,
		start_date,
		end_date,
		is_trx_created,
		quantity,
        created_by,
        created_at,
        updated_by,
        updated_at
	) values (
		$1, TO_DATE($2, 'DD/MM/YYYY'), $3, $4, $5, $6, TO_DATE($7, 'DD/MM/YYYY'), TO_DATE($8, 'DD/MM/YYYY'), $9, $10
		,$11, $12, $13, $14
	) RETURNING id`

	err := cmd.
		QueryRow(
			ctx,
			query,
			topup.CustomerId,
			topup.TransferDate,
			topup.TrxIdFrom,
			topup.TrxIdTo,
			topup.MerchantTo,
			topup.ChannelTo,
			topup.StartDate,
			topup.EndDate,
			topup.IsTrxCreated,
			topup.Quantity,
			topup.CreatedBy,
			time.Unix(0, topup.CreatedAt*int64(time.Millisecond)),
			topup.UpdatedBy,
			time.Unix(0, topup.UpdatedAt*int64(time.Millisecond)),
		).
		Scan(&id)

	log.Println("ALL TOP UP", topup)

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TransferBalanceRepositoryPostgre.create",
				"query": query,
			}).
			Error(err)

		return 0, pgxerror.FromPgxError(err, "", "TransferBalanceRepositoryPostgre.Create")
	}

	return id, nil
}

func (c *TransferBalanceRepositoryPostgre) Update(ctx context.Context, id int64, topup entity.TransferBalance, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update transfer_balance
	set
		customer_id = $1,
		transfer_date = $2,
        trx_id_from = $3,
        trx_id_to = $4,
                         "merchant_to" = $5,
                         "channel_to" = $6,
                         start_date = TO_DATE($7, 'DD/MM/YYYY'),
                         end_date = TO_DATE($8, 'DD/MM/YYYY'),
                         is_trx_created = $9,
                         quantity = $10,
                         updated_by = $11,
                         updated_at = $12
	where
		id = $13`

	_, err := cmd.Exec(
		ctx,
		query,
		topup.CustomerId,
		topup.TransferDate,
		topup.TrxIdFrom,
		topup.TrxIdTo,
		topup.MerchantTo,
		topup.ChannelTo,
		topup.StartDate,
		topup.EndDate,
		topup.IsTrxCreated,
		topup.Quantity,
		topup.UpdatedBy,
		topup.UpdatedAt,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "TransferBalanceRepositoryPostgre.Update")
	}

	return nil
}

func (c *TransferBalanceRepositoryPostgre) Update2(ctx context.Context, id int64, topup entity.TransferBalance, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update transfer_balance
	set
		so_number = $1,
		customer_id = $2,
        merchant_id = $3,
        channel_id = $4,
                         "start_date" = $5,
                         "end_date" = $6,
                         duration = $7,
                         billing = $8,
                         item_id = $9,
                         balance = $10,
                         rate = $11,
                         prepaid = $13,
                         quotation_id = $14,
                         void_date = $15,
                         amount = $16,
                         updated_by = $17,
                         updated_at = $18
	where
		id = $12`

	_, err := cmd.Exec(
		ctx,
		query,
		topup.CustomerId,
		topup.MerchantTo,
		topup.ChannelTo,
		topup.StartDate,
		topup.EndDate,
		topup.Quantity,
		topup.TransferDate,
		topup.IsTrxCreated,
		topup.TrxIdFrom,
		topup.TrxIdTo,
		topup.CreatedBy,
		topup.CreatedAt,
		topup.UpdatedBy,
		topup.UpdatedAt,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "TransferBalanceRepositoryPostgre.Update")
	}

	return nil
}

func (c *TransferBalanceRepositoryPostgre) Delete(ctx context.Context, id int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from transfer_balances where id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "TransferBalanceRepositoryPostgre.Delete")
	}

	return nil
}
