package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	if cust.CRMLeadID == "" {
		entityStatus = "6"
		// recordType = "customer"
	} else {
		entityStatus = "13"
		// recordType = "lead"
	}

	log.Println("Before entityStatus log statement")
	log.Println("entityStatus", entityStatus)
	log.Println("After entityStatus log statement")

	crdCustParam = credential.CustomerParam{
		Recordtype:                     "customer",
		Customform:                     "2",
		EntityID:                       cust.EnterprisePrivyID,
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		Comments:                       "",
		Email:                          cust.Email,
		EntityStatus:                   "13",
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

	var entityStatus string

	if cust.CRMLeadID == "" {
		entityStatus = "6"
	} else {
		entityStatus = "13"
	}

	fmt.Println("========= entityStatus ========", entityStatus)

	crdCustParam := credential.LeadParam{
		Recordtype:                     "lead",
		Customform:                     "2",
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		Comments:                       "",
		Email:                          cust.Email,
		EntityStatus:                   "6",
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

	customFields := map[string]interface{}{
		"Enterprise ID": cust.EnterprisePrivyID,
	}

	requestData := map[string]interface{}{
		"first_name":        cust.FirstName,
		"last_name":         cust.LastName,
		"email":             cust.Email,
		"organization_name": cust.CustomerName,
		"phone":             cust.PhoneNo,
		"custom_fields":     customFields,
	}

	jsonData, err := json.Marshal(requestData)

	if err != nil {
		// Handle error
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomer",
				"param": requestData,
			}).
			Error(err)
	}
	client := &http.Client{}
	// leadResp, err := http.Post("http://apergu.tech:9002/api/v1/zendesk/lead/on-create", "application/json", bytes.NewBuffer(requestData))
	leadResp, err := http.NewRequest("POST", "https://api.getbase.com/v2/deals", bytes.NewBuffer(jsonData))
	leadResp.Header.Set("Authorization", "26bed09778079a78eb96acb73feb1cb2d9b36267e992caa12b0d960c8f760e2c")
	leadResp.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(leadResp)

	log.Println("response", resp)
	log.Println("err", err)
	if err != nil {
		// Handle error
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomer",
				"param": requestData,
			}).
			Error(err)

		return 0, nil, err
	}

	defer resp.Body.Close()

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

	var entityStatus string

	if cust.CRMLeadID == "" {
		entityStatus = "6"
	} else {
		entityStatus = "13"
	}

	crdCustParam := credential.LeadParam{
		Recordtype:                     "lead",
		Customform:                     "2",
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		Comments:                       "",
		Email:                          cust.Email,
		EntityStatus:                   entityStatus,
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

	var entityStatus string

	if cust.CRMLeadID == "" {
		entityStatus = "6"
	} else {
		entityStatus = "13"
	}

	crdCustParam := credential.CustomerParam{
		Recordtype:                     "lead",
		Customform:                     "2",
		EntityID:                       cust.CRMLeadID,
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		Comments:                       "",
		Email:                          cust.Email,
		EntityStatus:                   entityStatus,
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

	var entityStatus string

	if cust.CRMLeadID == "" {
		entityStatus = "6"
	} else {
		entityStatus = "13"
	}

	crdCustParam := credential.CustomerParam{
		Recordtype:                     "lead",
		Customform:                     "2",
		EntityID:                       cust.CRMLeadID,
		IsPerson:                       "F",
		CompanyName:                    cust.CustomerName,
		Comments:                       "",
		Email:                          cust.Email,
		EntityStatus:                   entityStatus,
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
