package main

import (
	"encoding/json"
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
)

var repo repositories.JournalPgRepository

func init() {
	persisterConfig := mock.NewPersisterPostgresqlConfig()
	pst := microservice.NewPersister(persisterConfig)
	repo = repositories.NewJournalPgRepository(pst)
}

func main() {
	testUpdate()
}

func testUpdate() {
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
				"debitamount": 1200,
				"creditamount": 0
			},
			{
				"id" : 2,
				"accountcode": "11",
				"accountname": "11",
				"debitamount": 0,
				"creditamount": 1200
			}
		],
		"amount": 1000,
		"accountdescription": "",
		"bookcode": ""
	}`

	doc := models.JournalPg{}
	err := json.Unmarshal([]byte(journal_json), &doc)

	if err != nil {
		fmt.Errorf("error %v", err)
	}
	err = repo.Update(doc.ShopID, doc.DocNo, doc)
	if err != nil {
		fmt.Errorf("error %v", err)
	}
}

func testGetJournal() {
	data, err := repo.Get("27dcEdktOoaSBYFmnN6G6ett4Jb", "JO-202206067CFB22")
	if err != nil {
		fmt.Printf("error %v", err)
	}

	fmt.Printf("%v", data)
}

func testCreateJournal() {

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

	doc := models.JournalPg{}
	err := json.Unmarshal([]byte(journal_json), &doc)

	if err != nil {
		fmt.Errorf("error %v", err)
	}
	err = repo.Create(doc)
	if err != nil {
		fmt.Errorf("error %v", err)
	}
}
