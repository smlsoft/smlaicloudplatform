package repositories_test

import (
	"encoding/json"
	"os"
	"smlcloudplatform/internal/vfgl/journal/models"
	"smlcloudplatform/internal/vfgl/journal/repositories"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/microservice"
	"testing"

	"github.com/stretchr/testify/assert"
)

var journal_json string = `{
	"id": "000000000000000000000000",
	"shopid": "27dcEdktOoaSBYFmnN6G6ett4Jb",
	"guidfixed": "2ABh7CJyA7RbeZ1WmdwXWvs0GQa",
	"parid": "0000000",
	"batchId": "",
	"docno": "JO-202206067CFB22",
	"docdate": "2022-06-06T04:11:28.56Z",
	"accountperiod": 1,
	"accountyear": 2022,
	"accountgroup": "1",
	"journaldetail": [
		{
			"accountcode": "11010",
			"accountname": "เงินสด - บัญชี 1 (เงินล้าน) ",
			"debitamount": 1000,
			"creditamount": 0
		},
		{
			"accountcode": "11",
			"accountname": "11",
			"debitamount": 0,
			"creditamount": 1000
		}
	],
	"amount": 1000,
	"accountdescription": "",
	"bookcode": ""
}`

var repo repositories.JournalPgRepository

func init() {
	persisterConfig := mock.NewPersisterPostgresqlConfig()
	pst := microservice.NewPersister(persisterConfig)
	repo = repositories.NewJournalPgRepository(pst)
	pst.AutoMigrate(
		models.JournalPg{},
		models.JournalDetailPg{},
	)
}

func TestJournalRepositoryRealDBCreate(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	assert.NotNil(t, repo, "Failed to Init Repo")
	doc := models.JournalPg{}
	err := json.Unmarshal([]byte(journal_json), &doc)

	assert.Nil(t, err, "Failed Unmarshal Json to JournalDoc Create1")
	assert.Equal(t, doc.DocNo, "JO-202206067CFB22", "Failed Doc No Not Match")

	err = repo.Create(doc)
	assert.Nil(t, err, "Failed Unmarshal Json to JournalDoc Create1")
}

func TestJournalRepositoryRealDBGetAssertErrNotFound(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	assert.NotNil(t, repo, "Failed to Init Repo")

	_, err := repo.Get("TESTSHOP", "DOC01")
	assert.Error(t, err, "Failed Get Journal Doc NotFound But No Error")
}

func TestCreateAndDeleteJournal(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	json_str := `{
		"id": "62cdc14ca3f6ef3ca30543e8",
		"shopid": "2BYWCndV194TYXVEO7NlRLuJYWY",
		"guidfixed": "2Br5noZ5LmgRuQLwYrpreUl2J9a",
		"batchId": "",
		"docno": "JO-20220713014532831-1",
		"docdate": "2022-07-12T18:45:32.068Z",
		"documentref": "",
		"accountperiod": 7,
		"accountyear": 2565,
		"accountgroup": "2",
		"amount": 200,
		"accountdescription": "ชำระค่าธรรมเนียมสมัครสมาชิก#6500000049",
		"bookcode": "2",
		"vats": [],
		"taxes": [],
		"parid": "",
		"journaldetail": [
		  {
			"accountcode": "11020",
			"accountname": "เงินสด - บัญชี 2",
			"debitamount": 200,
			"creditamount": 0
		  },
		  {
			"accountcode": "43010",
			"accountname": "รายได้ - ค่าธรรมเนียม-แรกเข้า",
			"debitamount": 0,
			"creditamount": 200
		  }
		]
	  }	
	`
	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	assert.NotNil(t, repo, "Failed to Init Repo")

	doc := models.JournalPg{}
	err := json.Unmarshal([]byte(json_str), &doc)
	assert.Nil(t, err, "Failed Unmarshal Json to JournalDoc Create1")

	err = repo.Create(doc)
	assert.Nil(t, err, "Failed Unmarshal Json to JournalDoc Create1")

	err = repo.Delete(doc.ShopID, doc.DocNo)
	assert.Nil(t, err, "Failed Unmarshal Json to JournalDoc Create1")

}
