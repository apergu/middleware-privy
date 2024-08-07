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

type ApplicationCommandUsecaseGeneral struct {
	applicationRepo  repository.ApplicationRepository
	customerRepo     repository.CustomerRepository
	applicationPrivy credential.Application
}

func NewApplicationCommandUsecaseGeneral(prop ApplicationUsecaseProperty) *ApplicationCommandUsecaseGeneral {
	return &ApplicationCommandUsecaseGeneral{
		applicationRepo:  prop.ApplicationRepo,
		applicationPrivy: prop.ApplicationPrivy,
		customerRepo:     prop.CustomerRepo,
	}
}

func (r *ApplicationCommandUsecaseGeneral) Create(ctx context.Context, application model.Application) (int64, interface{}, error) {
	tx, err := r.applicationRepo.BeginTx(ctx)
	log.Println("customerRepo", r.customerRepo)
	if err != nil {
		return 0, nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			r.applicationRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.applicationRepo.RollbackTx(ctx, tx)
		} else {
			err = r.applicationRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	tmNow := time.Now().UnixNano() / 1000000

	// respCust, _ := r.applicationRepo.FindByName(ctx, merchant.MerchantName, tx)

	// if respCust.MerchantName != "" {
	// 	return 0, nil, rapperror.ErrConflict(
	// 		"",
	// 		"Merchant with name "+merchant.MerchantName+" already exist",
	// 		"MerchantCommandUsecaseGeneral.Create",
	// 		nil,
	// 	)
	// }

	// defer func() {
	// 	if p := recover(); p != nil {
	// 		r.applicationRepo.RollbackTx(ctx, tx)
	// 		panic(p)
	// 	} else if err != nil {
	// 		log.Println("Rolling back transaction due to error:", err)
	// 		r.applicationRepo.RollbackTx(ctx, tx)
	// 	} else {
	// 		err = r.applicationRepo.CommitTx(ctx, tx)
	// 		if err != nil {
	// 			log.Println("Error committing transaction:", err)
	// 		}
	// 	}
	// }()

	respEnterprise, _ := r.customerRepo.FindByEnterprisePrivyID(ctx, application.EnterpriseID, tx)

	if respEnterprise.EnterprisePrivyID == "" {
		return 0, nil, rapperror.ErrUnprocessableEntity(
			"",
			"Enterprise with ID "+application.EnterpriseID+" is Not Found",
			"ApplicationCommandUsecaseGeneral.Create",
			nil,
		)
	}

	defer func() {
		if p := recover(); p != nil {
			r.applicationRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.applicationRepo.RollbackTx(ctx, tx)
		} else {
			err = r.applicationRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	respCust2, _ := r.applicationRepo.FindByApplicationID(ctx, application.ApplicationID, tx)

	if respCust2.ApplicationID != "" && respCust2.ApplicationID != "000" {
		return 0, nil, rapperror.ErrConflict(
			"",
			"Application with Application ID "+application.ApplicationID+" already exist",
			"ApplicationCommandUsecaseGeneral.Create",
			nil,
		)
	}

	defer func() {
		if p := recover(); p != nil {
			r.applicationRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.applicationRepo.RollbackTx(ctx, tx)
		} else {
			err = r.applicationRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	insertApplication := entity.Application{
		CustomerID:      application.CustomerID,
		EnterpriseID:    application.EnterpriseID,
		ApplicationID:   application.ApplicationID,
		ApplicationName: application.ApplicationName,
		Address:         application.Address,
		Email:           application.Email,
		PhoneNo:         application.PhoneNo,
		State:           application.State,
		City:            application.City,
		ZipCode:         application.ZipCode,
		CreatedBy:       application.CreatedBy,
		CreatedAt:       tmNow,
		UpdatedBy:       application.CreatedBy,
		UpdatedAt:       tmNow,
		ApplicationCode: application.ApplicationCode,
	}

	applicationId, err := r.applicationRepo.Create(ctx, insertApplication, tx)
	print("applicationId", applicationId)
	if err != nil {
		r.applicationRepo.CommitTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationCommandUsecaseGeneral.Create",
				"src":   "customerRepo.Create",
				"param": insertApplication,
			}).
			Error(err)

		data := map[string]interface{}{
			"application": insertApplication,
			"message":     err.Error(),
		}

		return 0, data, nil
	}

	defer func() {
		if p := recover(); p != nil {
			r.applicationRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.applicationRepo.RollbackTx(ctx, tx)
		} else {
			err = r.applicationRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	// find customer by application.EnterpriseID
	customer_filter := repository.CustomerFilter{
		EnterprisePrivyID: &application.EnterpriseID,
	}
	customers, _ := r.customerRepo.Find(ctx, customer_filter, 1, 0, nil)

	defer func() {
		if p := recover(); p != nil {
			r.applicationRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.applicationRepo.CommitTx(ctx, tx)
		} else {
			err = r.applicationRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	var customer entity.Customer
	if len(customers) > 0 {
		customer = customers[0]
	}

	// if customer.CustomerInternalID == 0 {

	// }

	// custrecordcustomer_name ambil dari customer

	privyParam := credential.ApplicationParam{
		RecordType:                     "customrecord_customer_hierarchy",
		CustRecordCustomerName:         customer.CustomerInternalID,
		CustRecordEnterpriseID:         application.EnterpriseID,
		CustRecordApplicationID:        application.ApplicationID,
		CustRecordPrivyCodeApplication: application.ApplicationCode,
		CustRecordApplicationName:      application.ApplicationName,
		CustRecordAddress:              application.Address,
		CustRecordEmail:                application.Email,
		CustRecordPhone:                application.PhoneNo,
		CustRecordState:                application.State,
		CustRecordCity:                 application.City,
		CustRecordZip:                  application.ZipCode,
		Method:                         "POST",
	}

	resp, err := r.applicationPrivy.CreateApplication(ctx, privyParam)
	if err != nil {
		r.applicationRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationCommandUsecaseGeneral.Create",
				"src":   "applicationPrivy.CreateApplication",
				"param": privyParam,
			}).
			Error(err)

		return 0, nil, err
	}

	insertApplication.ApplicationInternalID = resp.Data.RecordID
	insertApplication.CustomerInternalID = customer.CustomerInternalID

	err = r.applicationRepo.Update(ctx, applicationId, insertApplication, tx)
	if err != nil {
		r.applicationRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationCommandUsecaseGeneral.Create",
				"src":   "customerRepo.Update",
				"param": insertApplication,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.applicationRepo.CommitTx(ctx, tx)
	if err != nil {
		r.applicationRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "MerchantCommandUsecaseGeneral.Create",
				"src": "customerRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"MerchantCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return applicationId, nil, nil
}

func (r *ApplicationCommandUsecaseGeneral) Update(ctx context.Context, id int64, application model.Application) (int64, interface{}, error) {
	tx, err := r.applicationRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	print("applicationsFIND2", tmNow)

	updatedApplication := entity.Application{
		ApplicationCode: application.ApplicationCode,
		ApplicationID:   application.ApplicationID,
		ApplicationName: application.ApplicationName,
		Address:         application.Address,
		Email:           application.Email,
		PhoneNo:         application.PhoneNo,
		State:           application.State,
		City:            application.City,
		ZipCode:         application.ZipCode,
		UpdatedBy:       application.CreatedBy,
	}

	err = r.applicationRepo.Update(ctx, id, updatedApplication, tx)
	if err != nil {
		r.applicationRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.applicationRepo.CommitTx(ctx, tx)
	if err != nil {
		r.applicationRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "ApplicationCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"ApplicationCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *ApplicationCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.applicationRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.applicationRepo.Delete(ctx, id, tx)
	if err != nil {
		r.applicationRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationCommandUsecaseGeneral.Delete",
				"src":   "custRepo.Delete",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.applicationRepo.CommitTx(ctx, tx)
	if err != nil {
		r.applicationRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "ApplicationCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"ApplicationCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
