package usecase

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/mohamadikbal/project-privy/pkg/credential"
	"gitlab.com/rteja-library3/rapperror"
)

type CustomerUsageCommandUsecaseGeneral struct {
	custUsageRepo repository.CustomerUsageCommandRepository
	customerPrivy credential.Customer
}

func NewCustomerUsageCommandUsecaseGeneral(prop CustomerUsageUsecaseProperty) *CustomerUsageCommandUsecaseGeneral {
	return &CustomerUsageCommandUsecaseGeneral{
		custUsageRepo: prop.CustomerUsageRepo,
		customerPrivy: prop.CustomerPrivy,
	}
}

func (r *CustomerUsageCommandUsecaseGeneral) Create(ctx context.Context, cust model.CustomerUsage) (int64, interface{}, error) {
	tx, err := r.custUsageRepo.BeginTx(ctx)
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

	custId, err := r.custUsageRepo.Create(ctx, insertCustomerUsage, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertCustomerUsage,
			}).
			Error(err)

		return 0, nil, err
	}

	crdCustParam := credential.CustomerParam{
		Recordtype:                     "customer",
		Customform:                     "2",
		EntityID:                       cust.CustomerID,
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		EntityStatus:                   cust.EntityStatus,
		Comments:                       "",
		URL:                            cust.URL,
		Email:                          cust.Email,
		Phone:                          cust.Phone,
		AltPhone:                       cust.AltPhone,
		Fax:                            cust.Fax,
		CustEntityPrivyCustomerBalance: int(cust.Balance),
		CustEntityPrivyCustomerUsage:   int(cust.Usage),
	}

	if cust.IsPerson {
		crdCustParam.IsPerson = "T"
	}

	_, err = r.customerPrivy.CreateCustomer(ctx, crdCustParam)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomer",
				"param": crdCustParam,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custUsageRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

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

	return custId, nil, nil
}

func (r *CustomerUsageCommandUsecaseGeneral) Update(ctx context.Context, id int64, cust model.CustomerUsage) (int64, interface{}, error) {
	tx, err := r.custUsageRepo.BeginTx(ctx)
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

	err = r.custUsageRepo.Update(ctx, id, updatedCustomerUsage, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custUsageRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

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
	tx, err := r.custUsageRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.custUsageRepo.Delete(ctx, id, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Delete",
				"src":   "custRepo.Delete",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custUsageRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

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
