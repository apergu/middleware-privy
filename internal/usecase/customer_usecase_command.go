package usecase

import (
	"context"
	"log"
	"time"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type CustomerCommandUsecaseGeneral struct {
	custRepo      repository.CustomerCommandRepository
	customerPrivy credential.Customer
}

func NewCustomerCommandUsecaseGeneral(prop CustomerUsecaseProperty) *CustomerCommandUsecaseGeneral {
	return &CustomerCommandUsecaseGeneral{
		custRepo:      prop.CustomerRepo,
		customerPrivy: prop.CustomerPrivy,
	}
}

func (r *CustomerCommandUsecaseGeneral) Create(ctx context.Context, cust model.Customer) (int64, interface{}, error) {

	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	insertCustomer := entity.Customer{
		CustomerID:        cust.EnterprisePrivyID,
		CustomerType:      cust.CustomerType,
		CustomerName:      cust.CustomerName,
		FirstName:         cust.FirstName,
		LastName:          cust.LastName,
		Email:             cust.Email,
		PhoneNo:           cust.PhoneNo,
		Address:           cust.Address,
		CRMLeadID:         cust.CRMLeadID,
		EnterprisePrivyID: cust.EnterprisePrivyID,
		NPWP:              cust.NPWP,
		Address1:          cust.Address1,
		State:             cust.State,
		City:              cust.City,
		CreatedBy:         cust.CreatedBy,
		CreatedAt:         tmNow,
		UpdatedBy:         cust.CreatedBy,
		UpdatedAt:         tmNow,
	}

	custId, err := r.custRepo.Create(ctx, insertCustomer, tx)
	log.Println("response", err)

	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertCustomer,
			}).
			Error(err)

		return 0, nil, err
	}

	crdCustParam := credential.CustomerParam{
		Recordtype:                     "customer",
		Customform:                     "2",
		EntityID:                       cust.EnterprisePrivyID,
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		Comments:                       "",
		Email:                          cust.Email,
		EntityStatus:                   cust.EntityStatus,
		URL:                            cust.URL,
		Phone:                          cust.PhoneNo,
		AltPhone:                       cust.AltPhone,
		Fax:                            cust.Fax,
		CustEntityPrivyCustomerBalance: cust.Balance,
		CustEntityPrivyCustomerUsage:   cust.Usage,
		EnterprisePrivyID:              cust.EnterprisePrivyID,
		NPWP:                           cust.NPWP,
		Address1:                       cust.Address1,
		State:                          cust.State,
		City:                           cust.City,
		ZipCode:                        cust.ZipCode,
		CompanyNameLong:                cust.CustomerName,
		CRMLeadID:                      cust.CRMLeadID,
		BankAccount:                    "103",
		AddressBook: credential.AddressBook{
			Addr1: cust.Address1,
			State: cust.State,
			City:  cust.City,
			Zip:   cust.ZipCode,
		},
	}

	privyResp, err := r.customerPrivy.CreateCustomer(ctx, crdCustParam)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomer",
				"param": crdCustParam,
			}).
			Error(err)

		return 0, nil, err
	}

	insertCustomer.CustomerInternalID = privyResp.Details.CustomerInternalID
	err = r.custRepo.Update(ctx, custId, insertCustomer, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": custId,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

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

func (r *CustomerCommandUsecaseGeneral) CreateLead2(ctx context.Context, cust model.Customer) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	insertCustomer := entity.Customer{
		CustomerID:        cust.CRMLeadID,
		CustomerType:      cust.CustomerType,
		CustomerName:      cust.CustomerName,
		FirstName:         cust.FirstName,
		LastName:          cust.LastName,
		Email:             cust.Email,
		PhoneNo:           cust.PhoneNo,
		Address:           cust.Address,
		CRMLeadID:         cust.CRMLeadID,
		EnterprisePrivyID: cust.EnterprisePrivyID,
		NPWP:              cust.NPWP,
		Address1:          cust.Address1,
		State:             cust.State,
		City:              cust.City,
		CreatedBy:         cust.CreatedBy,
		CreatedAt:         tmNow,
		UpdatedBy:         cust.CreatedBy,
		UpdatedAt:         tmNow,
	}

	custId, err := r.custRepo.CreateLead(ctx, insertCustomer, tx)
	log.Println("response", err)

	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertCustomer,
			}).
			Error(err)

		return 0, nil, err
	}

	crdCustParam := credential.CustomerParam{
		Recordtype:                     "lead",
		Customform:                     "2",
		EntityID:                       cust.CRMLeadID,
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		Comments:                       "",
		Email:                          cust.Email,
		EntityStatus:                   cust.EntityStatus,
		URL:                            cust.URL,
		Phone:                          cust.PhoneNo,
		AltPhone:                       cust.AltPhone,
		Fax:                            cust.Fax,
		CustEntityPrivyCustomerBalance: cust.Balance,
		CustEntityPrivyCustomerUsage:   cust.Usage,
		EnterprisePrivyID:              cust.EnterprisePrivyID,
		NPWP:                           cust.NPWP,
		Address1:                       cust.Address1,
		State:                          cust.State,
		City:                           cust.City,
		ZipCode:                        cust.ZipCode,
		CompanyNameLong:                cust.CustomerName,
		CRMLeadID:                      cust.CRMLeadID,
		BankAccount:                    "103",
		AddressBook: credential.AddressBook{
			Addr1: cust.Address1,
			State: cust.State,
			City:  cust.City,
			Zip:   cust.ZipCode,
		},
	}

	privyResp, err := r.customerPrivy.CreateLead(ctx, crdCustParam)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomer",
				"param": crdCustParam,
			}).
			Error(err)

		return 0, nil, err
	}

	insertCustomer.CustomerInternalID = privyResp.Details.CustomerInternalID
	err = r.custRepo.Update(ctx, custId, insertCustomer, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": custId,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

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

func (r *CustomerCommandUsecaseGeneral) CreateLead(ctx context.Context, cust model.Lead) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	insertCustomer := entity.Customer{
		CustomerID:        cust.CRMLeadID,
		CustomerType:      cust.CustomerType,
		CustomerName:      cust.CustomerName,
		FirstName:         cust.FirstName,
		LastName:          cust.LastName,
		Email:             cust.Email,
		PhoneNo:           cust.PhoneNo,
		Address:           cust.Address,
		CRMLeadID:         cust.CRMLeadID,
		EnterprisePrivyID: cust.EnterprisePrivyID,
		NPWP:              cust.NPWP,
		Address1:          cust.Address1,
		State:             cust.State,
		City:              cust.City,
		CreatedBy:         cust.CreatedBy,
		CreatedAt:         tmNow,
		UpdatedBy:         cust.CreatedBy,
		UpdatedAt:         tmNow,
	}

	custId, err := r.custRepo.CreateLead(ctx, insertCustomer, tx)
	log.Println("response", err)

	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertCustomer,
			}).
			Error(err)

		return 0, nil, err
	}

	crdCustParam := credential.CustomerParam{
		Recordtype:                     "lead",
		Customform:                     "2",
		EntityID:                       cust.CRMLeadID,
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		Comments:                       "",
		Email:                          cust.Email,
		EntityStatus:                   cust.EntityStatus,
		URL:                            cust.URL,
		Phone:                          cust.PhoneNo,
		AltPhone:                       cust.AltPhone,
		Fax:                            cust.Fax,
		CustEntityPrivyCustomerBalance: cust.Balance,
		CustEntityPrivyCustomerUsage:   cust.Usage,
		EnterprisePrivyID:              cust.EnterprisePrivyID,
		NPWP:                           cust.NPWP,
		Address1:                       cust.Address1,
		State:                          cust.State,
		City:                           cust.City,
		ZipCode:                        cust.ZipCode,
		CompanyNameLong:                cust.CustomerName,
		CRMLeadID:                      cust.CRMLeadID,
		BankAccount:                    "103",
		AddressBook: credential.AddressBook{
			Addr1: cust.Address1,
			State: cust.State,
			City:  cust.City,
			Zip:   cust.ZipCode,
		},
	}

	privyResp, err := r.customerPrivy.CreateLead(ctx, crdCustParam)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomer",
				"param": crdCustParam,
			}).
			Error(err)

		return 0, nil, err
	}

	insertCustomer.CustomerInternalID = privyResp.Details.CustomerInternalID
	err = r.custRepo.Update(ctx, custId, insertCustomer, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": custId,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

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

func (r *CustomerCommandUsecaseGeneral) UpdateLead(ctx context.Context, id string, cust model.Lead) (any, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	updatedCustomer := entity.Customer{
		CustomerID:        cust.EnterprisePrivyID,
		CustomerType:      cust.CustomerType,
		CustomerName:      cust.CustomerName,
		FirstName:         cust.FirstName,
		LastName:          cust.LastName,
		Email:             cust.Email,
		PhoneNo:           cust.PhoneNo,
		Address:           cust.Address,
		CRMLeadID:         cust.CRMLeadID,
		EnterprisePrivyID: cust.EnterprisePrivyID,
		NPWP:              cust.NPWP,
		Address1:          cust.Address1,
		State:             cust.State,
		City:              cust.City,
		UpdatedBy:         cust.CreatedBy,
		UpdatedAt:         tmNow,
	}

	err = r.custRepo.UpdateLead(ctx, id, updatedCustomer, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	crdCustParam := credential.CustomerParam{
		Recordtype:                     "lead",
		Customform:                     "2",
		EntityID:                       cust.CRMLeadID,
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		Comments:                       "",
		Email:                          cust.Email,
		EntityStatus:                   cust.EntityStatus,
		URL:                            cust.URL,
		Phone:                          cust.PhoneNo,
		AltPhone:                       cust.AltPhone,
		Fax:                            cust.Fax,
		CustEntityPrivyCustomerBalance: cust.Balance,
		CustEntityPrivyCustomerUsage:   cust.Usage,
		EnterprisePrivyID:              cust.EnterprisePrivyID,
		NPWP:                           cust.NPWP,
		Address1:                       cust.Address1,
		State:                          cust.State,
		City:                           cust.City,
		ZipCode:                        cust.ZipCode,
		CompanyNameLong:                cust.CustomerName,
		CRMLeadID:                      cust.CRMLeadID,
		BankAccount:                    "103",
		AddressBook: credential.AddressBook{
			Addr1: cust.Address1,
			State: cust.State,
			City:  cust.City,
			Zip:   cust.ZipCode,
		},
	}

	privyResp, err := r.customerPrivy.UpdateLead(ctx, crdCustParam)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomer",
				"param": crdCustParam,
			}).
			Error(err)

		return 0, nil, err
	}

	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": privyResp,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

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

func (r *CustomerCommandUsecaseGeneral) UpdateLead2(ctx context.Context, id int64, cust model.Lead) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	updatedCustomer := entity.Customer{
		CustomerID:        cust.CRMLeadID,
		CustomerType:      cust.CustomerType,
		CustomerName:      cust.CustomerName,
		FirstName:         cust.FirstName,
		LastName:          cust.LastName,
		Email:             cust.Email,
		PhoneNo:           cust.PhoneNo,
		Address:           cust.Address,
		CRMLeadID:         cust.CRMLeadID,
		EnterprisePrivyID: cust.CRMLeadID,
		NPWP:              cust.NPWP,
		Address1:          cust.Address1,
		State:             cust.State,
		City:              cust.City,
		UpdatedBy:         cust.CreatedBy,
		UpdatedAt:         tmNow,
	}

	err = r.custRepo.Update(ctx, id, updatedCustomer, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	crdCustParam := credential.CustomerParam{
		Recordtype:                     "lead",
		Customform:                     "2",
		EntityID:                       cust.CRMLeadID,
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		Comments:                       "",
		Email:                          cust.Email,
		EntityStatus:                   cust.EntityStatus,
		URL:                            cust.URL,
		Phone:                          cust.PhoneNo,
		AltPhone:                       cust.AltPhone,
		Fax:                            cust.Fax,
		CustEntityPrivyCustomerBalance: cust.Balance,
		CustEntityPrivyCustomerUsage:   cust.Usage,
		EnterprisePrivyID:              cust.EnterprisePrivyID,
		NPWP:                           cust.NPWP,
		Address1:                       cust.Address1,
		State:                          cust.State,
		City:                           cust.City,
		ZipCode:                        cust.ZipCode,
		CompanyNameLong:                cust.CustomerName,
		CRMLeadID:                      cust.CRMLeadID,
		BankAccount:                    "103",
		AddressBook: credential.AddressBook{
			Addr1: cust.Address1,
			State: cust.State,
			City:  cust.City,
			Zip:   cust.ZipCode,
		},
	}

	//privyResp, err := r.customerPrivy.UpdateLead(ctx, crdCustParam)

	privyResp, err := r.customerPrivy.UpdateLead(ctx, crdCustParam)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomer",
				"param": crdCustParam,
			}).
			Error(err)

		return 0, nil, err
	}

	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": privyResp,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

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

func (r *CustomerCommandUsecaseGeneral) Update(ctx context.Context, id int64, cust model.Customer) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	updatedCustomer := entity.Customer{
		CustomerID:        cust.CRMLeadID,
		CustomerType:      cust.CustomerType,
		CustomerName:      cust.CustomerName,
		FirstName:         cust.FirstName,
		LastName:          cust.LastName,
		Email:             cust.Email,
		PhoneNo:           cust.PhoneNo,
		Address:           cust.Address,
		CRMLeadID:         cust.CRMLeadID,
		EnterprisePrivyID: cust.CRMLeadID,
		NPWP:              cust.NPWP,
		Address1:          cust.Address1,
		State:             cust.State,
		City:              cust.City,
		UpdatedBy:         cust.CreatedBy,
		UpdatedAt:         tmNow,
	}

	err = r.custRepo.Update(ctx, id, updatedCustomer, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Update",
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

func (r *CustomerCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.custRepo.Delete(ctx, id, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Delete",
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
				"at":  "CustomerCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"CustomerCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
