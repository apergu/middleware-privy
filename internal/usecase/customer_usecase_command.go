package usecase

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
	merchantRepo  repository.MerchantQueryRepository
	channelRepo   repository.ChannelQueryRepository
	customerPrivy credential.Customer
	channelPrivy  credential.Channel
	merchantPrivy credential.Merchant
}

func NewCustomerCommandUsecaseGeneral(prop CustomerUsecaseProperty) *CustomerCommandUsecaseGeneral {
	return &CustomerCommandUsecaseGeneral{
		custRepo:      prop.CustomerRepo,
		customerPrivy: prop.CustomerPrivy,
		merchantPrivy: prop.MerchantPrivy,
		channelPrivy:  prop.ChannelPrivy,
		merchantRepo:  prop.MerchantRepo,
		channelRepo:   prop.ChannelRepo,
	}
}

func (r *CustomerCommandUsecaseGeneral) Create(ctx context.Context, cust model.Customer) (int64, interface{}, error) {

	log.Println("BEFORE ERRROR ", r.merchantRepo)
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	respCRMLead, _ := r.custRepo.FindByCRMLeadId(ctx, cust.CRMLeadID, tx)

	if respCRMLead.CRMLeadID != "" {
		return 0, nil, rapperror.ErrConflict(
			"",
			"Customer with CRM Lead ID "+cust.CRMLeadID+" already Won",
			"CustomerCommandUsecaseGeneral.Create",
			nil,
		)
	}

	defer func() {
		if p := recover(); p != nil {
			r.custRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.custRepo.RollbackTx(ctx, tx)
		} else {
			err = r.custRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	fmt.Println("respCust FIND NAME")
	respCust, _ := r.custRepo.FindByName(ctx, cust.CustomerName, tx)

	if respCust.CustomerName != "" {
		return 0, nil, rapperror.ErrConflict(
			"",
			"Customer with name "+cust.CustomerName+" already exist",
			"CustomerCommandUsecaseGeneral.Create",
			nil,
		)
	}

	defer func() {
		if p := recover(); p != nil {
			r.custRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.custRepo.RollbackTx(ctx, tx)
		} else {
			err = r.custRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	fmt.Println("respCust FIND NAME")
	respCust2, _ := r.custRepo.FindByEnterprisePrivyID(ctx, cust.EnterprisePrivyID, tx)

	if respCust2.EnterprisePrivyID != "" {
		return 0, nil, rapperror.ErrConflict(
			"",
			"Customer with enterprise privy id "+cust.EnterprisePrivyID+" already exist",
			"CustomerCommandUsecaseGeneral.Create",
			nil,
		)
	}

	defer func() {
		if p := recover(); p != nil {
			r.custRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.custRepo.RollbackTx(ctx, tx)
		} else {
			err = r.custRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()
	log.Println("merchant CUST TEST", cust.EnterprisePrivyID)

	merchant, err := r.merchantRepo.FindByEnterprisePrivyID(ctx, cust.EnterprisePrivyID, nil)

	log.Println("merchant", err)
	log.Println("merchant2", merchant)

	defer func() {
		if p := recover(); p != nil {
			r.custRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.custRepo.RollbackTx(ctx, tx)
		} else {
			err = r.custRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	if merchant.MerchantID != "" {
		payloadMerchant := credential.MerchantParam{
			RecordType:                  "customrecord_customer_hierarchy",
			CustRecordCustomerName:      0,
			CustRecordEnterpriseID:      merchant.EnterpriseID,
			CustRecordMerchantID:        merchant.MerchantID,
			CustRecordPrivyCodeMerchant: merchant.MerchantCode,
			CustRecordMerchantName:      merchant.MerchantName,
			CustRecordAddress:           merchant.Address,
			CustRecordEmail:             merchant.Email,
			CustRecordPhone:             merchant.PhoneNo,
			CustRecordState:             merchant.State,
			CustRecordCity:              merchant.City,
			CustRecordZip:               merchant.ZipCode,
			Method:                      "POST",
		}

		r.merchantPrivy.CreateMerchant(ctx, payloadMerchant)

		channel, _ := r.channelRepo.FindByMerchantID(ctx, merchant.MerchantID, nil)

		defer func() {
			if p := recover(); p != nil {
				r.custRepo.RollbackTx(ctx, tx)
				panic(p)
			} else if err != nil {
				log.Println("Rolling back transaction due to error:", err)
				r.custRepo.RollbackTx(ctx, tx)
			} else {
				err = r.custRepo.CommitTx(ctx, tx)
				if err != nil {
					log.Println("Error committing transaction:", err)
				}
			}
		}()

		payloadChannel := credential.ChannelParam{
			RecordType:                 "customrecord_customer_hierarchy",
			CustRecordCustomerName:     strconv.Itoa(int(merchant.CustomerInternalID)),
			CustRecordEnterpriseID:     merchant.EnterpriseID,
			CustRecordChannelID:        channel.ChannelID,
			CustRecordMerchantID:       merchant.MerchantID,
			CustRecordPrivyCodeChannel: channel.ChannelCode,
			CustRecordChannelName:      channel.ChannelName,
			CustRecordAddress:          channel.Address,
			CustRecordEmail:            channel.Email,
			CustRecordPhone:            channel.PhoneNo,
			CustRecordState:            channel.State,
			CustRecordCity:             channel.City,
			CustRecordZip:              channel.ZipCode,
			Method:                     channel.Method,
		}

		r.channelPrivy.CreateChannel(ctx, payloadChannel)
	}

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

	defer func() {
		if p := recover(); p != nil {
			r.custRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.custRepo.RollbackTx(ctx, tx)
		} else {
			err = r.custRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()
	log.Println("response", err)
	log.Println("BEFORE ERRROR ")

	defer func() {
		if p := recover(); p != nil {
			r.custRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.custRepo.RollbackTx(ctx, tx)
		} else {
			err = r.custRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

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
		SubIndustry:                    cust.SubIndustry,
		CRMLeadID:                      cust.CRMLeadID,
		BankAccount:                    "103",
		AddressBook: credential.AddressBook{
			Addr1: cust.Address,
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
		// SubIndustry:                    cust.SubIndustry,
		CRMLeadID:   cust.CRMLeadID,
		BankAccount: "103",
		AddressBook: credential.AddressBook{
			Addr1: cust.Address1,
			State: cust.State,
			City:  cust.City,
			Zip:   cust.ZipCode,
		},
	}

	// customFields := map[string]interface{}{
	// 	"Enterprise ID": cust.EnterprisePrivyID,
	// }

	// requestData := map[string]interface{}{
	// 	"first_name":        cust.FirstName,
	// 	"last_name":         cust.LastName,
	// 	"email":             cust.Email,
	// 	"organization_name": cust.CustomerName,
	// 	"phone":             cust.PhoneNo,
	// 	"custom_fields":     customFields,
	// }

	// dataReq := map[string]interface{}{
	// 	"data": requestData,
	// }

	// jsonData, err := json.Marshal(dataReq)

	// if err != nil {
	// 	// Handle error
	// 	r.custRepo.RollbackTx(ctx, tx)

	// 	logrus.
	// 		WithFields(logrus.Fields{
	// 			"at":    "CustomerCommandUsecaseGeneral.Create",
	// 			"src":   "customerPrivy.CreateCustomer",
	// 			"param": requestData,
	// 		}).
	// 		Error(err)
	// }
	// client := &http.Client{}
	// leadResp, _ := http.NewRequest("POST", "https://api.getbase.com/v2/leads", bytes.NewBuffer(jsonData))
	// leadResp.Header.Set("Authorization", "Bearer 26bed09778079a78eb96acb73feb1cb2d9b36267e992caa12b0d960c8f760e2c")
	// leadResp.Header.Set("Content-Type", "application/json")

	// resp, err := client.Do(leadResp)

	// log.Println("response", resp)
	// log.Println("err", err)
	// if err != nil {
	// 	// Handle error
	// 	r.custRepo.RollbackTx(ctx, tx)

	// 	logrus.
	// 		WithFields(logrus.Fields{
	// 			"at":    "CustomerCommandUsecaseGeneral.Create",
	// 			"src":   "customerPrivy.CreateCustomer",
	// 			"param": requestData,
	// 		}).
	// 		Error(err)

	// 	return 0, nil, err
	// }

	// defer resp.Body.Close()

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

func (r *CustomerCommandUsecaseGeneral) CreateLeadZD(ctx context.Context, cust model.Customer) (int64, interface{}, error) {
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

	// var entityStatus string

	// if cust.CRMLeadID == "" {
	// 	entityStatus = "6"
	// } else {
	// 	entityStatus = "13"
	// }

	// fmt.Println("========= entityStatus ========", entityStatus)

	// crdCustParam := credential.LeadParam{
	// 	Recordtype:                     "lead",
	// 	Customform:                     "2",
	// 	IsPerson:                       "F",
	// 	CompanyName:                    cust.CustomerName,
	// 	Comments:                       "",
	// 	Email:                          cust.Email,
	// 	EntityStatus:                   "6",
	// 	URL:                            cust.URL,
	// 	Phone:                          cust.PhoneNo,
	// 	AltPhone:                       cust.AltPhone,
	// 	Fax:                            cust.Fax,
	// 	CustEntityPrivyCustomerBalance: cust.Balance,
	// 	CustEntityPrivyCustomerUsage:   cust.Usage,
	// 	EnterprisePrivyID:              cust.EnterprisePrivyID,
	// 	NPWP:                           cust.NPWP,
	// 	Address1:                       cust.Address1,
	// 	State:                          cust.State,
	// 	City:                           cust.City,
	// 	ZipCode:                        cust.ZipCode,
	// 	CompanyNameLong:                cust.CustomerName,
	// 	CRMLeadID:                      cust.CRMLeadID,
	// 	BankAccount:                    "103",
	// 	AddressBook: credential.AddressBook{
	// 		Addr1: cust.Address1,
	// 		State: cust.State,
	// 		City:  cust.City,
	// 		Zip:   cust.ZipCode,
	// 	},
	// }

	// // customFields := map[string]interface{}{
	// // 	"Enterprise ID": cust.EnterprisePrivyID,
	// // }

	// // requestData := map[string]interface{}{
	// // 	"first_name":        cust.FirstName,
	// // 	"last_name":         cust.LastName,
	// // 	"email":             cust.Email,
	// // 	"organization_name": cust.CustomerName,
	// // 	"phone":             cust.PhoneNo,
	// // 	"custom_fields":     customFields,
	// // }

	// // dataReq := map[string]interface{}{
	// // 	"data": requestData,
	// // }

	// // jsonData, err := json.Marshal(dataReq)

	// // if err != nil {
	// // 	// Handle error
	// // 	r.custRepo.RollbackTx(ctx, tx)

	// // 	logrus.
	// // 		WithFields(logrus.Fields{
	// // 			"at":    "CustomerCommandUsecaseGeneral.Create",
	// // 			"src":   "customerPrivy.CreateCustomer",
	// // 			"param": requestData,
	// // 		}).
	// // 		Error(err)
	// // }
	// // client := &http.Client{}
	// // leadResp, _ := http.NewRequest("POST", "https://api.getbase.com/v2/leads", bytes.NewBuffer(jsonData))
	// // leadResp.Header.Set("Authorization", "Bearer 26bed09778079a78eb96acb73feb1cb2d9b36267e992caa12b0d960c8f760e2c")
	// // leadResp.Header.Set("Content-Type", "application/json")

	// // resp, err := client.Do(leadResp)

	// // log.Println("response", resp)
	// // log.Println("err", err)
	// // if err != nil {
	// // 	// Handle error
	// // 	r.custRepo.RollbackTx(ctx, tx)

	// // 	logrus.
	// // 		WithFields(logrus.Fields{
	// // 			"at":    "CustomerCommandUsecaseGeneral.Create",
	// // 			"src":   "customerPrivy.CreateCustomer",
	// // 			"param": requestData,
	// // 		}).
	// // 		Error(err)

	// // 	return 0, nil, err
	// // }

	// // defer resp.Body.Close()

	// privyResp, err := r.customerPrivy.CreateLead(ctx, crdCustParam)
	// if err != nil {
	// 	r.custRepo.RollbackTx(ctx, tx)

	// 	logrus.
	// 		WithFields(logrus.Fields{
	// 			"at":    "CustomerCommandUsecaseGeneral.Create",
	// 			"src":   "customerPrivy.CreateCustomer",
	// 			"param": crdCustParam,
	// 		}).
	// 		Error(err)

	// 	return 0, nil, err
	// }

	// insertCustomer.CustomerInternalID = privyResp.Details.CustomerInternalID
	// err = r.custRepo.Update(ctx, custId, insertCustomer, tx)
	// if err != nil {
	// 	r.custRepo.RollbackTx(ctx, tx)

	// 	logrus.
	// 		WithFields(logrus.Fields{
	// 			"at":    "CustomerCommandUsecaseGeneral.Create",
	// 			"src":   "custRepo.Update",
	// 			"param": custId,
	// 		}).
	// 		Error(err)

	// 	return 0, nil, err
	// }

	err = r.custRepo.CommitTx(ctx, tx)

	fmt.Println("========= ERROR ========", err)
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
		// SubIndustry:                    cust.SubIndustry,
		CRMLeadID:   cust.CRMLeadID,
		BankAccount: "103",
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
