package usecase

import (
	"errors"
	common "smlcloudplatform/internal/models"
	transmodels "smlcloudplatform/internal/transaction/models"
	"smlcloudplatform/internal/transaction/payment/models"
)

func ParseTransactionToPayment(other transmodels.TransactionMessageQueue) (models.TransactionPayment, error) {

	doc := models.TransactionPayment{}

	doc.ShopID = other.ShopID
	doc.DocNo = other.DocNo
	doc.DocDate = other.DocDatetime
	doc.GuidRef = other.GuidRef
	doc.TransFlag = int8(other.TransFlag)
	doc.DocType = other.DocType
	doc.InquiryType = other.InquiryType
	doc.IsCancel = other.IsCancel
	doc.PayCashAmount = other.PayCashAmount

	tempBranch := common.JSONB{}
	err := tempBranch.Scan(other.Branch)

	if err != nil {
		return models.TransactionPayment{}, err
	}

	doc.Branch = tempBranch

	tempTransFlag, err := transFlagToCalcFlag(other.TransFlag)

	if err != nil {
		return models.TransactionPayment{}, err
	}

	doc.CalcFlag = tempTransFlag

	doc.PayCashChange = other.PayCashChange
	doc.SumQRCode = other.SumQRCode
	doc.SumCreditCard = other.SumCreditCard
	doc.SumMoneyTransfer = other.SumMoneyTransfer
	doc.SumCheque = other.SumCheque
	doc.SumCoupon = other.SumCoupon
	doc.TotalAmount = other.TotalAmount
	doc.RoundAmount = other.RoundAmount
	doc.SumCredit = other.SumCredit

	return doc, nil
}

func transFlagToCalcFlag(transFlag int) (int8, error) {

	paid := []int{16, 44, 50} // รับชำระ 239->50
	pay := []int{12, 48, 51}  // จ่ายชำระ 19->51

	if contains(paid, transFlag) {
		return 1, nil
	} else if contains(pay, transFlag) {
		return -1, nil
	}

	return 0, errors.New("invalid transflag")
}

func contains(arr []int, val int) bool {

	for _, v := range arr {
		if v == val {
			return true
		}
	}

	return false
}
