package usecase_test

import (
	"smlcloudplatform/internal/models"
	transmodels "smlcloudplatform/internal/transaction/models"
	payment_models "smlcloudplatform/internal/transaction/payment/models"
	"smlcloudplatform/internal/transaction/payment/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTransactionToPayment(t *testing.T) {

	langCode := "en"
	langName := "BranchTest"

	branch := []models.NameX{
		{
			Code: &langCode,
			Name: &langName,
		},
	}

	testCases := []struct {
		name     string
		doc      transmodels.TransactionMessageQueue
		expected payment_models.TransactionPayment
	}{
		{
			name: "Parse Transaction to Payment",
			doc: transmodels.TransactionMessageQueue{
				Transaction: transmodels.Transaction{
					TransactionHeader: transmodels.TransactionHeader{
						DocNo:            "doc1",
						PaymentDetailRaw: "{\"cashamount\":1.0,\"cashamounttext\":\"\",\"paymentcreditcards\":[],\"paymenttransfers\":[]}",
						TransFlag:        50,
						Branch:           transmodels.TransactionBranch{Code: "0001", Names: &branch},
					},
				},
			},
			expected: payment_models.TransactionPayment{
				BranchCode:  "0001",
				BranchNames: branch,
				TransFlag:   50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paymentDoc, err := usecase.ParseTransactionToPayment(tc.doc)

			assert.NoError(t, err, "Error should be nil")
			assert.Equal(t, tc.expected.BranchCode, paymentDoc.BranchCode, "Branch code invalid")
			assert.Equal(t, tc.expected.BranchNames, paymentDoc.BranchNames, "Branch name invalid")
			assert.Equal(t, "en", *paymentDoc.BranchNames[0].Code, "Branch name code should be en")
			assert.Equal(t, int8(50), paymentDoc.TransFlag, "Trans Flag should be 50")
		})
	}

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

	branch := transmodels.JSONBTransactionBranch{}
	err := branch.Scan([]byte(`{"guidfixed":"","code":"0001","names":[{"code":"en","name":"Test","isauto":false,"isdelete":false}]}`))
	// err = branch.Scan(tempBranch)

	assert.Nil(t, err, "Error should be nil")

}
