package repositories_test

import (
	"encoding/json"
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
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
