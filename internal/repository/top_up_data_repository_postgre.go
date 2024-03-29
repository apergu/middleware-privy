package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

type TopUpDataRepositoryPostgre struct {
	pool *pgxpool.Pool
}

func NewTopUpDataRepositoryPostgre(pool *pgxpool.Pool) *TopUpDataRepositoryPostgre {
	return &TopUpDataRepositoryPostgre{
		pool: pool,
	}
}

func (c *TopUpDataRepositoryPostgre) query(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) ([]entity.TopUpData, error) {
	rows, err := cmd.Query(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataRepositoryPostgre.query",
				"src":   "cmd.Query",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when query",
			"TopUpDataRepositoryPostgre.query.Query",
			nil,
		)
	}
	defer rows.Close()

	result := make([]entity.TopUpData, 0)
	for rows.Next() {
		data := entity.TopUpData{}

		err := rows.Scan(
			&data.ID,
			&data.MerchantID,
			&data.TransactionID,
			&data.EnterpriseID,
			&data.EnterpriseName,
			&data.OriginalServiceID,
			&data.ServiceID,
			&data.ServiceName,
			&data.Quantity,
			&data.TransactionDate,
			&data.MerchantCode,
			&data.ChannelID,
			&data.ChannelCode,
			&data.CustomerInternalID,
			&data.MerchantInternalID,
			&data.ChannelInternalID,
			&data.TransactionType,
			&data.TopupID,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "TopUpDataRepositoryPostgre.query",
					"src":   "rows.Scan",
					"query": query,
					"args":  args,
				}).
				Error(err)

			return nil, pgxerror.FromPgxError(
				err,
				"Something went wrong when scan",
				"TopUpDataRepositoryPostgre.query.Scan",
			)
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *TopUpDataRepositoryPostgre) queryOne(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (entity.TopUpData, error) {
	data := entity.TopUpData{}

	err := cmd.
		QueryRow(ctx, query, args...).
		Scan(
			&data.ID,
			&data.MerchantID,
			&data.TransactionID,
			&data.EnterpriseID,
			&data.EnterpriseName,
			&data.OriginalServiceID,
			&data.ServiceID,
			&data.ServiceName,
			&data.Quantity,
			&data.TransactionDate,
			&data.MerchantCode,
			&data.ChannelID,
			&data.ChannelCode,
			&data.CustomerInternalID,
			&data.MerchantInternalID,
			&data.ChannelInternalID,
			&data.TransactionType,
			&data.TopupID,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataRepositoryPostgre.queryOne",
				"src":   "rows.Scan",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return entity.TopUpData{}, pgxerror.FromPgxError(
			err,
			"Something went wrong when scan",
			"TopUpDataRepositoryPostgre.queryOne.Scan",
		)
	}

	return data, nil
}

func (c *TopUpDataRepositoryPostgre) buildFilter(filter TopUpDataFilter) (string, []interface{}) {
	condBuilder := &strings.Builder{}
	conds := make([]string, 0, 4) // set for 2 capacity is posible max filter
	condArgs := make([]interface{}, 0, 4)

	if len(conds) > 0 {
		condBuilder.WriteString("where ")
		condBuilder.WriteString(strings.Join(conds, " and "))
	}

	return condBuilder.String(), condArgs
}

func (c *TopUpDataRepositoryPostgre) buildSort(sort string) string {
	switch sort {
	case "newest":
		return `order by top_up_datas.created_at desc`
	}

	return `order by top_up_datas.updated_at desc`
}

func (c *TopUpDataRepositoryPostgre) Find(ctx context.Context, filter TopUpDataFilter, limit, skip int64, tx pgx.Tx) ([]entity.TopUpData, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	var limits, skips string
	cond, args := c.buildFilter(filter)

	order := c.buildSort(filter.Sort)

	query := `select
		top_up_datas."id",
		top_up_datas.merchant_id,
		top_up_datas.transaction_id,
		top_up_datas.enterprise_id,
		top_up_datas."enterprise_name",
		top_up_datas.original_service_id,
		top_up_datas."service_id",
		top_up_datas.service_name,
		top_up_datas.quantity,
		top_up_datas.transaction_date,
		top_up_datas.merchant_code,
		top_up_datas.channel_id,
		top_up_datas.channel_code,
		top_up_datas.customer_internalid,
		top_up_datas.merchant_internalid,
		top_up_datas.channel_internalid,
		top_up_datas.transaction_type,
		top_up_datas.topup_id,
		top_up_datas.created_by,
		top_up_datas.created_at,
		top_up_datas.updated_by,
		top_up_datas.updated_at
	from
		top_up_datas
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

func (c *TopUpDataRepositoryPostgre) Count(ctx context.Context, filter TopUpDataFilter, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	cond, args := c.buildFilter(filter)

	query := `select
		count(top_up_datas.id)
	from
		top_up_datas
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
			"TopUpDataRepositoryPostgre.Count.Scan",
		)
	}

	return data, nil
}

func (c *TopUpDataRepositoryPostgre) FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.TopUpData, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `select
		top_up_datas."id",
		top_up_datas.merchant_id,
		top_up_datas.transaction_id,
		top_up_datas.enterprise_id,
		top_up_datas."enterprise_name",
		top_up_datas.original_service_id,
		top_up_datas."service_id",
		top_up_datas.service_name,
		top_up_datas.quantity,
		top_up_datas.transaction_date,
		top_up_datas.merchant_code,
		top_up_datas.channel_id,
		top_up_datas.channel_code,
		top_up_datas.customer_internalid,
		top_up_datas.merchant_internalid,
		top_up_datas.channel_internalid,
		top_up_datas.transaction_type,
		top_up_datas.topup_id,
		top_up_datas.created_by,
		top_up_datas.created_at,
		top_up_datas.updated_by,
		top_up_datas.updated_at
	from
		top_up_datas
	where
		top_up_datas.id = $1
	limit 1`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *TopUpDataRepositoryPostgre) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return c.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (c *TopUpDataRepositoryPostgre) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (c *TopUpDataRepositoryPostgre) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (c *TopUpDataRepositoryPostgre) FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.TopUpData, error) {
	if tx == nil {
		return entity.TopUpData{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"TopUpDataRepositoryPostgre.FindOneByIdForUpdate",
			nil,
		)
	}
	var cmd sqlcommand.Command = tx

	query := `select
		top_up_datas."id",
		top_up_datas.merchant_id,
		top_up_datas.transaction_id,
		top_up_datas.enterprise_id,
		top_up_datas."enterprise_name",
		top_up_datas.original_service_id,
		top_up_datas."service_id",
		top_up_datas.service_name,
		top_up_datas.quantity,
		top_up_datas.transaction_date,
		top_up_datas.merchant_code,
		top_up_datas.channel_id,
		top_up_datas.channel_code,
		top_up_datas.customer_internalid,
		top_up_datas.merchant_internalid,
		top_up_datas.channel_internalid,
		top_up_datas.transaction_type,
		top_up_datas.topup_id,
		top_up_datas.created_by,
		top_up_datas.created_at,
		top_up_datas.updated_by,
		top_up_datas.updated_at
	from
		top_up_datas
	where
		top_up_datas.id = $1
	limit 1
	FOR UPDATE`

	return c.queryOne(ctx, cmd, query, id)
}

func (c *TopUpDataRepositoryPostgre) Create(ctx context.Context, topup entity.TopUp, tx pgx.Tx) (int64, error) {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	topUpUUid := uuid.New().String()

	var id int64
	query := `insert into top_up_data (
		top_up_id,
		so_number,
		customer_id,
        merchant_id,
        channel_id,
                         start_date,
                         end_date,
                         duration,
                         billing,
                         item_id,
                         balance,
                         rate,
                         prepaid,
                         quotation_id,
                         void_date,
                         amount,
                         created_by,
                         created_at,
                         updated_by,
                         updated_at
	) values (
		$1, $2, $3, $4, $5, TO_DATE($6, 'DD/MM/YYYY'), TO_DATE($7, 'DD/MM/YYYY'), $8, $9, $10
		,$11, $12, $13, $14, TO_DATE($15, 'DD/MM/YYYY'), $16, $17, $18, $19, $20
	) RETURNING id`

	err := cmd.
		QueryRow(
			ctx,
			query,
			topUpUUid,
			topup.SoNo,
			topup.CustomerId,
			topup.MerchantId,
			topup.ChannelId,
			topup.StartDate,
			topup.EndDate,
			topup.Duration,
			topup.Billing,
			topup.ItemId,
			topup.QtyBalance,
			topup.Rate,
			topup.Prepaid,
			topup.QuotationId,
			topup.VoidDate,
			topup.Amount,
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
				"at":    "TopUpDataRepositoryPostgre.create",
				"query": query,
			}).
			Error(err)

		return 0, pgxerror.FromPgxError(err, "", "TopUpDataRepositoryPostgre.Create")
	}

	return id, nil
}

func (c *TopUpDataRepositoryPostgre) Update(ctx context.Context, id int64, topup entity.TopUp, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update top_up_data
	set
		so_number = $1,
		customer_id = $2,
        merchant_id = $3,
        channel_id = $4,
                         "start_date" = TO_DATE($5, 'DD/MM/YYYY'),
                         "end_date" = TO_DATE($6, 'DD/MM/YYYY'),
                         duration = $7,
                         billing = $8,
                         item_id = $9,
                         balance = $10,
                         rate = $11,
                         prepaid = $13,
                         quotation_id = $14,
                         void_date = TO_DATE($15, 'DD/MM/YYYY'),
                         amount = $16,
                         updated_by = $17,
                         updated_at = $18
	where
		id = $12`

	_, err := cmd.Exec(
		ctx,
		query,
		topup.SoNo,
		topup.CustomerId,
		topup.MerchantId,
		topup.ChannelId,
		topup.StartDate,
		topup.EndDate,
		topup.Duration,
		topup.Billing,
		topup.ItemId,
		topup.QtyBalance,
		topup.Rate,
		id,
		topup.Prepaid,
		topup.QuotationId,
		topup.VoidDate,
		topup.Amount,
		topup.UpdatedBy,
		topup.UpdatedAt,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "TopUpDataRepositoryPostgre.Update")
	}

	return nil
}

func (c *TopUpDataRepositoryPostgre) Update2(ctx context.Context, id int64, topup entity.TopUpData, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := `update top_up_data
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
		topup.MerchantID,
		topup.TransactionID,
		topup.EnterpriseID,
		topup.EnterpriseName,
		topup.OriginalServiceID,
		topup.ServiceID,
		topup.ServiceName,
		topup.Quantity,
		topup.TransactionDate,
		topup.UpdatedBy,
		topup.UpdatedAt,
		id,
		topup.MerchantCode,
		topup.ChannelID,
		topup.ChannelCode,
		topup.CustomerInternalID,
		topup.MerchantInternalID,
		topup.ChannelInternalID,
		topup.TransactionType,
		topup.TopupID,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "TopUpDataRepositoryPostgre.Update")
	}

	return nil
}

func (c *TopUpDataRepositoryPostgre) Delete(ctx context.Context, id int64, tx pgx.Tx) error {
	var cmd sqlcommand.Command = c.pool
	if tx != nil {
		cmd = tx
	}

	query := "delete from top_up_datas where id = $1"
	_, err := cmd.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return pgxerror.FromPgxError(err, "", "TopUpDataRepositoryPostgre.Delete")
	}

	return nil
}
