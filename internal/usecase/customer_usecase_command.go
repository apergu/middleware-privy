package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"
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

	// var lastId int64

	// idLast, err := r.custRepo.GetLast(ctx, tx)
	// if err != nil {
	// 	r.custRepo.RollbackTx(ctx, tx)

	// 	logrus.
	// 		WithFields(logrus.Fields{
	// 			"at":  "CustomerCommandUsecaseGeneral.Create",
	// 			"src": "custRepo.GetLast",
	// 		}).
	// 		Error(err)

	// }

	// id := idLast.ID + 1
	// id = strconv.FormatInt(id, 6)
	insertCustomer := InsertCustomer(cust)

	if cust.CustomerType != nil {
		insertCustomer.CustomerType = *cust.CustomerType
	}

	custId, err := r.custRepo.Create(ctx, insertCustomer, tx)
	log.Println("response", err)
	log.Println("BEFORE ERRROR ")

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

	log.Println("After ERRROR ")

	var entityStatus string
	// var recordType string

	var crdCustParam credential.CustomerParam

	if cust.CRMLeadID != nil && *cust.CRMLeadID == "" {
		entityStatus = "6"
		// recordType = "customer"
	} else {
		entityStatus = "13"
		// recordType = "lead"
	}

	log.Println("Before entityStatus log statement")
	log.Println("entityStatus", entityStatus)
	log.Println("After entityStatus log statement")

	crdCustParam = CrdInsertCustomer(cust, "13")
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

	insertCustomer := InsertCustomer(cust)

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

	var entityStatus string

	if cust.CRMLeadID != nil && *cust.CRMLeadID == "" {
		entityStatus = "6"
	} else {
		entityStatus = "13"
	}

	fmt.Println("========= entityStatus ========", entityStatus)

	crdCustParam := CrdInsertCustomer(cust, "6")
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

func (r *CustomerCommandUsecaseGeneral) CreateLead(ctx context.Context, cust model.Customer) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	insertCustomer := InsertCustomer(cust)
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

	var entityStatus string

	if *cust.CRMLeadID == "" {
		entityStatus = "6"
	} else {
		entityStatus = "13"
	}

	crdCustParam := CrdInsertCustomer(cust, entityStatus)

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

func (r *CustomerCommandUsecaseGeneral) UpdateLead(ctx context.Context, id string, cust model.Customer) (any, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	updatedCustomer := InsertCustomer(cust)

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

	var entityStatus string

	if *cust.CRMLeadID == "" {
		entityStatus = "6"
	} else {
		entityStatus = "13"
	}

	crdCustParam := CrdInsertCustomer(cust, entityStatus)

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

func (r *CustomerCommandUsecaseGeneral) UpdateLead2(ctx context.Context, id int64, cust model.Customer) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}
	updatedCustomer := InsertCustomer(cust)
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

	var entityStatus string

	if *cust.CRMLeadID == "" {
		entityStatus = "6"
	} else {
		entityStatus = "13"
	}

	crdCustParam := CrdInsertCustomer(cust, entityStatus)

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
		CustomerID:        *cust.CRMLeadID,
		CustomerType:      *cust.CustomerType,
		CustomerName:      *cust.CustomerName,
		FirstName:         *cust.FirstName,
		LastName:          *cust.LastName,
		Email:             *cust.Email,
		PhoneNo:           *cust.PhoneNo,
		Address:           *cust.Address,
		CRMLeadID:         *cust.CRMLeadID,
		EnterprisePrivyID: *cust.CRMLeadID,
		NPWP:              *cust.NPWP,
		Address1:          *cust.Address1,
		State:             *cust.State,
		City:              *cust.City,
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

func InsertCustomer(cust model.Customer) entity.Customer {
	// Assuming tmNow is meant to be the current time
	tmNow := time.Now().UnixNano() / 999999

	// Create and return a new Customer entity
	return entity.Customer{
		CustomerID:        safeDeref(cust.EnterprisePrivyID),
		CustomerName:      safeDeref(cust.CustomerName),
		FirstName:         safeDeref(cust.FirstName),
		LastName:          safeDeref(cust.LastName),
		Email:             safeDeref(cust.Email),
		PhoneNo:           safeDeref(cust.PhoneNo),
		Address:           safeDeref(cust.Address),
		CRMLeadID:         safeDeref(cust.CRMLeadID),
		EnterprisePrivyID: safeDeref(cust.EnterprisePrivyID),
		NPWP:              safeDeref(cust.NPWP),
		Address1:          safeDeref(cust.Address1),
		State:             safeDeref(cust.State),
		City:              safeDeref(cust.City),
		CreatedBy:         cust.CreatedBy,
		CreatedAt:         tmNow,
		UpdatedBy:         cust.CreatedBy,
		UpdatedAt:         tmNow,
	}
}

func CrdInsertCustomer(cust model.Customer, entitySts string) credential.CustomerParam {
	return credential.CustomerParam{
		Recordtype:                     "customer",
		Customform:                     "2",
		EntityID:                       safeDeref(cust.EnterprisePrivyID),
		IsPerson:                       "F",
		CompanyName:                    safeDeref(cust.CustomerName),
		Comments:                       "",
		Email:                          safeDeref(cust.Email),
		EntityStatus:                   entitySts,
		URL:                            safeDeref(cust.URL),
		Phone:                          safeDeref(cust.PhoneNo),
		AltPhone:                       cust.AltPhone,
		Fax:                            cust.Fax,
		CustEntityPrivyCustomerBalance: cust.Balance,
		CustEntityPrivyCustomerUsage:   cust.Usage,
		EnterprisePrivyID:              safeDeref(cust.EnterprisePrivyID),
		NPWP:                           safeDeref(cust.NPWP),
		Address1:                       safeDeref(cust.Address1),
		State:                          safeDeref(cust.State),
		City:                           safeDeref(cust.City),
		ZipCode:                        safeDeref(cust.ZipCode),
		CompanyNameLong:                safeDeref(cust.CustomerName),
		CRMLeadID:                      safeDeref(cust.CRMLeadID),
		BankAccount:                    "103",
		AddressBook: credential.AddressBook{
			Addr1: safeDeref(cust.Address1),
			State: safeDeref(cust.State),
			City:  safeDeref(cust.City),
			Zip:   safeDeref(cust.ZipCode),
		},
	}
}

func safeDeref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
