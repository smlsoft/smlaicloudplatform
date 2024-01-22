package usecase

import (
	"encoding/json"
	transmodels "smlcloudplatform/internal/transaction/models"
	"smlcloudplatform/internal/transaction/paymentdetail/models"
)

func ParseTransactionToPaymentDetail(other transmodels.TransactionMessageQueue) ([]models.TransactionPaymentDetail, error) {

	tempDetails := []models.TransactionPaymentDetail{}
	rawDetail := other.PaymentDetailRaw

	if rawDetail == "" || rawDetail == "[]" || rawDetail == "null" {
		return []models.TransactionPaymentDetail{}, nil
	}

	err := json.Unmarshal([]byte(rawDetail), &tempDetails)

	if err != nil {
		return []models.TransactionPaymentDetail{}, err
	}

	for i := range tempDetails {
		tempDetails[i].ShopID = other.ShopID
		tempDetails[i].DocNo = other.DocNo

		// switch field
		tempDetails[i].PaymentType = tempDetails[i].TransFlag
		tempDetails[i].TransFlag = other.TransFlag
	}

	return tempDetails, nil
}
