package usecase

import (
	"context"
	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/mohamadikbal/project-privy/pkg/credential"
	"gitlab.com/rteja-library3/rapperror"
	"log"
	//"time"
)

type LeadCommandUsecaseGeneral struct {
	leadRepo  repository.LeadCommandRepository
	leadPrivy credential.Lead
}

func NewLeadCommandUsecaseGeneral(prop LeadUsecaseProperty) *LeadCommandUsecaseGeneral {
	return &LeadCommandUsecaseGeneral{
		leadRepo:  prop.LeadRepo,
		leadPrivy: prop.LeadPrivy,
	}
}

func (r *LeadCommandUsecaseGeneral) Create(ctx context.Context, cust model.Leads) (int64, interface{}, error) {
	tx, err := r.leadRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	//tmNow := time.Now().UnixNano() / 1000000

	insertCustomer := entity.Leads{
		//RecordType:   "lead",
		CustomerID:   cust.CustomerID,
		EnterpriseId: cust.EnterpriseId,
		IsPerson:     "F",
		CompanyName:  cust.CompanyName,
		Email:        cust.Email,
		Phone:        cust.Phone,
		Fax:          cust.Fax,
		NPWP:         cust.NPWP,
		CRMLeadID:    cust.CRMLeadID,
		BankAccount:  "103",
	}

	custId, err := r.leadRepo.Create(ctx, insertCustomer, tx)
	log.Println("response", err)

	if err != nil {
		r.leadRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertCustomer,
			}).
			Error(err)

		return 0, nil, err
	}

	crdCustParam := credential.LeadParam{
		RecordType:   "lead",
		EnterpriseId: cust.EnterpriseId,
		IsPerson:     "F",
		CompanyName:  cust.CompanyName,
		Email:        cust.Email,
		Phone:        cust.Phone,
		Fax:          cust.Fax,
		NPWP:         cust.NPWP,
		CRMLeadID:    cust.CRMLeadID,
		BankAccount:  "103",
	}

	privyResp, err := r.leadPrivy.CreateLead(ctx, crdCustParam)
	if err != nil {
		r.leadRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomer",
				"param": crdCustParam,
			}).
			Error(err)

		return 0, nil, err
	}

	insertCustomer.CustomerID = privyResp.Details.Customerid
	err = r.leadRepo.Update(ctx, custId, insertCustomer, tx)
	if err != nil {
		r.leadRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": custId,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.leadRepo.CommitTx(ctx, tx)
	if err != nil {
		r.leadRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "CustomerCommandUsecaseGeneral.Create",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"CustomerCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return custId, nil, nil
}

//
//
//func (r *LeadCommandUsecaseGeneral) Create(ctx context.Context, cust model.Leads) (int64, interface{}, error) {
//	tx, err := r.leadRepo.BeginTx(ctx)
//	if err != nil {
//		return 0, nil, err
//	}
//
//	//tmNow := time.Now().UnixNano() / 1000000
//
//	insertCustomer := entity.Leads{
//RecordType:   "lead",
//		EnterpriseId: cust.EnterpriseId,
//		IsPerson:     "F",
//		CompanyName:  cust.CompanyName,
//		Email:        cust.Email,
//		Phone:        cust.Phone,
//		Fax:          cust.Fax,
//		NPWP:         cust.NPWP,
//		CRMLeadID:    cust.CRMLeadID,
//		BankAccount:  "103",
//	}
//
//	custId, err := r.leadRepo.Create(ctx, insertCustomer, tx)
//	log.Println("response", err)
//
//	if err != nil {
//		r.leadRepo.RollbackTx(ctx, tx)
//
//		logrus.
//			WithFields(logrus.Fields{
//				"at":    "CustomerCommandUsecaseGeneral.Create",
//				"src":   "custRepo.Create",
//				"param": insertCustomer,
//			}).
//			Error(err)
//
//		return 0, nil, err
//	}
//
//	crdCustParam := credential.LeadParam{
//		RecordType:   "lead",
//		EnterpriseId: cust.EnterpriseId,
//		IsPerson:     "F",
//		CompanyName:  cust.CompanyName,
//		Email:        cust.Email,
//		Phone:        cust.Phone,
//		Fax:          cust.Fax,
//		NPWP:         cust.NPWP,
//		CRMLeadID:    cust.CRMLeadID,
//		BankAccount:  "103",
//	}
//
//	privyResp, err := r.leadPrivy.CreateLead(ctx, crdCustParam)
//	if err != nil {
//		r.leadRepo.RollbackTx(ctx, tx)
//
//		logrus.
//			WithFields(logrus.Fields{
//				"at":    "CustomerCommandUsecaseGeneral.Create",
//				"src":   "customerPrivy.CreateCustomer",
//				"param": crdCustParam,
//			}).
//			Error(err)
//
//		return 0, nil, err
//	}
//
//	insertCustomer.CustomerID = privyResp.Details.Customerid
//	//err = r.leadRepo.Update(ctx, custId, insertCustomer, tx)
//	if err != nil {
//		r.leadRepo.RollbackTx(ctx, tx)
//
//		logrus.
//			WithFields(logrus.Fields{
//				"at":    "CustomerCommandUsecaseGeneral.Create",
//				"src":   "custRepo.Update",
//				"param": custId,
//			}).
//			Error(err)
//
//		return 0, nil, err
//	}
//
//	err = r.leadRepo.CommitTx(ctx, tx)
//	if err != nil {
//		r.leadRepo.RollbackTx(ctx, tx)
//
//		logrus.
//			WithFields(logrus.Fields{
//				"at":  "CustomerCommandUsecaseGeneral.Create",
//				"src": "custRepo.CommitTx",
//			}).
//			Error(err)
//
//		return 0, nil, rapperror.ErrInternalServerError(
//			"",
//			"Something went wrong when commit",
//			"CustomerCommandUsecaseGeneral.Create",
//			nil,
//		)
//	}
//
//	return custId, nil, nil
//}

func (r *LeadCommandUsecaseGeneral) Update(ctx context.Context, id int64, cust model.Leads) (int64, interface{}, error) {
	tx, err := r.leadRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	//tmNow := time.Now().UnixNano() / 1000000

	updatedCustomer := entity.Leads{
		EnterpriseId: cust.EnterpriseId,
		IsPerson:     "F",
		CompanyName:  cust.CompanyName,
		Email:        cust.Email,
		Phone:        cust.Phone,
		Fax:          cust.Fax,
		NPWP:         cust.NPWP,
		CRMLeadID:    cust.CRMLeadID,
		BankAccount:  "103",
	}

	err = r.leadRepo.Update(ctx, id, updatedCustomer, tx)
	if err != nil {
		r.leadRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.leadRepo.CommitTx(ctx, tx)
	if err != nil {
		r.leadRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "CustomerCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"CustomerCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *LeadCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.leadRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.leadRepo.Delete(ctx, id, tx)
	if err != nil {
		r.leadRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "LeadCommandUsecaseGeneral.Delete",
				"src":   "custRepo.Delete",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.leadRepo.CommitTx(ctx, tx)
	if err != nil {
		r.leadRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "LeadCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"LeadCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
