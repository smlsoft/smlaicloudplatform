package usecase_test

import (
	"smlcloudplatform/internal/models"
	transmodels "smlcloudplatform/internal/transaction/models"
	"smlcloudplatform/internal/transaction/payment/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTransactionToPayment(t *testing.T) {

	doc := transmodels.TransactionMessageQueue{}
	doc.PaymentDetailRaw = "{\"cashamount\":0.0,\"cashamounttext\":\"\",\"paymentcreditcards\":[],\"paymenttransfers\":[]}"
	doc.TransFlag = 1

	nameCode := "en"
	name := "Test"

	doc.Branch.Code = "0001"
	doc.Branch.Names = &[]models.NameX{
		{
			Code: &nameCode,
			Name: &name,
		},
	}

	_, err := usecase.ParseTransactionToPayment(doc)

	assert.Nil(t, err, "Error should be nil")

}

func TestScanBytes(t *testing.T) {

	branchData := transmodels.TransactionBranch{}
	branchData.Code = "0001"

	nameCode := "en"
	name := "Test"
	branchData.Names = &[]models.NameX{
		{
			Code: &nameCode,
			Name: &name,
		},
	}

	// tempBranch, err := json.Marshal(branchData)

	// if err != nil {
	// 	assert.Nil(t, err, "Error should be nil")
	// }

	branch := transmodels.JSONBTransactionBranch{}
	// err := branch.Scan([]byte(`[{"code":"en","name":"Test"}]`))
	err := branch.Scan([]byte(`{"guidfixed":"","code":"0001","names":[{"code":"en","name":"Test","isauto":false,"isdelete":false}]}`))
	// err = branch.Scan(tempBranch)

	assert.Nil(t, err, "Error should be nil")

}
