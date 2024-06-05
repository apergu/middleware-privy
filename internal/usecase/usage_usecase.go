package usecase

import (
	"context"
	"middleware/internal/model"
	"middleware/pkg/credential"
	"strings"
)

type UsageUsecaseGeneral struct {
	customerUsagePrivy credential.CustomerUsage
}

func NewUsageUsecaseGeneral(prop CustomerUsageUsecaseProperty) *UsageUsecaseGeneral {
	return &UsageUsecaseGeneral{
		customerUsagePrivy: prop.CustomerPrivy,
	}
}

func (r *UsageUsecaseGeneral) Create(ctx context.Context, cust model.UsageModel) (*credential.CustomerUsageResponse, error) {

	custUsage := strings.Split(cust.TransactionID, "/")
	custPrivyUsgParam := credential.CustomerUsageParam{
		RecordType:                          "customrecord_privy_integrasi_usage",
		CustrecordPrivyServiceIntegrasi:     cust.ServiceID,
		CustrecordPrivyQuantityIntegrasi:    int64(cust.Qty),
		CustrecordPrivyUsageDateIntegrasi:   cust.TransactionDate,
		CustrecordEnterpriseeID:             custUsage[0],
		CustrecordPrivyMerchantNameIntgrasi: custUsage[1],
		CustrecordPrivyChannelNameIntgrasi:  custUsage[2],
	}

	res, err := r.customerUsagePrivy.CreateCustomerUsage(ctx, custPrivyUsgParam)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
