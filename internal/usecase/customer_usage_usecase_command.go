package usecase

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/rteja-library3/rapperror"
)

type CustomerUsageCommandUsecaseGeneral struct {
	custRepo repository.CustomerUsageCommandRepository
}

func NewCustomerUsageCommandUsecaseGeneral(prop CustomerUsageUsecaseProperty) *CustomerUsageCommandUsecaseGeneral {
	return &CustomerUsageCommandUsecaseGeneral{
		custRepo: prop.CustomerUsageRepo,
	}
}

func (r *CustomerUsageCommandUsecaseGeneral) Create(ctx context.Context, cust model.CustomerUsage) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000
	transAt := cust.TransactionAt.UnixNano() / 1000000

	insertCustomerUsage := entity.CustomerUsage{
		CustomerID:    cust.CustomerID,
		CustomerName:  cust.CustomerName,
		ProductID:     cust.ProductID,
		ProductName:   cust.ProductName,
		TransactionAt: transAt,
		Balance:       cust.Balance,
		BalanceAmount: cust.BalanceAmount,
		Usage:         cust.Usage,
		UsageAmount:   cust.UsageAmount,
		CreatedBy:     cust.CreatedBy,
		CreatedAt:     tmNow,
		UpdatedBy:     cust.CreatedBy,
		UpdatedAt:     tmNow,
	}

	roleId, err := r.custRepo.Create(ctx, insertCustomerUsage, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertCustomerUsage,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "CustomerUsageCommandUsecaseGeneral.Create",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"CustomerUsageCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return roleId, nil, nil
}

func (r *CustomerUsageCommandUsecaseGeneral) Update(ctx context.Context, id int64, cust model.CustomerUsage) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000
	transAt := cust.TransactionAt.UnixNano() / 1000000

	updatedCustomerUsage := entity.CustomerUsage{
		CustomerID:    cust.CustomerID,
		CustomerName:  cust.CustomerName,
		ProductID:     cust.ProductID,
		ProductName:   cust.ProductName,
		TransactionAt: transAt,
		Balance:       cust.Balance,
		BalanceAmount: cust.BalanceAmount,
		Usage:         cust.Usage,
		UsageAmount:   cust.UsageAmount,
		UpdatedBy:     cust.CreatedBy,
		UpdatedAt:     tmNow,
	}

	err = r.custRepo.Update(ctx, id, updatedCustomerUsage, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "CustomerUsageCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"CustomerUsageCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *CustomerUsageCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.custRepo.Delete(ctx, id, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Delete",
				"src":   "custRepo.Delete",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "CustomerUsageCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"CustomerUsageCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
