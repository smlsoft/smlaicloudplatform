package main

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/models"
	"sync"
	"time"

	"github.com/apex/log"
)

func main() {

	pstMongoCfg := mock.NewPersisterMongoConfig()
	pstPgCfg := mock.NewPersisterPostgresqlConfig()
	mqConfig := mock.NewMqConfig()

	logctx := log.WithFields(log.Fields{
		"name": "runx App ",
	})

	pstPg := microservice.NewPersister(pstPgCfg)
	pst := microservice.NewPersisterMongo(pstMongoCfg)
	prod := microservice.NewProducer(mqConfig.URI(), logctx)

	invRepo := inventory.NewInventoryRepository(pst)
	invPgRepo := inventory.NewInventoryIndexPGRepository(pstPg)
	invMqRepo := inventory.NewInventoryMQRepository(prod)
	invService := inventory.NewInventoryService(invRepo, invPgRepo, invMqRepo)

	shopID := "shoptest"
	authUser := "devtest"
	inv := models.Inventory{}
	inv.ItemSku = "itemSku"
	inv.CategoryGuid = "cate001"
	inv.Barcode = "bx001"
	inv.Price = 5
	inv.MemberPrice = 3
	inv.Name1 = "item test 1"
	inv.Description1 = "desc 1"

	opt := models.Option{}
	opt.Code = "opt1"
	opt.Required = true
	opt.SelectMode = "SIGLE"
	opt.MaxSelect = 1
	opt.Name1 = "opt name"

	inv.Options = &[]models.Option{
		opt,
	}

	syncMutex := sync.Mutex{}

	count := 0
	totalGen := 2
	// totalGen := 1000
	c := make(chan bool, 20)

	start := time.Now()
	for i := 1; i <= totalGen; i++ {
		go func(i int, c chan bool) {
			// fmt.Println(time.Now().Second())
			inv.ItemSku = fmt.Sprintf("itemSku%d", i)
			inv.Barcode = fmt.Sprintf("bx%d", i)

			// startMongo := time.Now()
			idx, guidx, err := invService.CreateInventory(shopID, authUser, inv)
			// fmt.Printf("mongo :: %s\n", time.Since(startMongo))
			if err != nil {
				println(err.Error())
			}

			docIdx := models.InventoryIndex{}
			docIdx.ID = idx
			docIdx.ShopID = shopID
			docIdx.GuidFixed = guidx

			// startPg := time.Now()
			err = invService.CreateIndex(docIdx)
			// fmt.Printf("pg :: %s\n", time.Since(startPg))

			if err != nil {
				println(err.Error())
			}
			syncMutex.Lock()
			// fmt.Printf("\r:: %d/%d", i, totalGen)
			count++
			syncMutex.Unlock()

			<-c
		}(i, c)

		c <- true

	}
	time.Sleep(time.Second * 1)
	print("\n\n")
	fmt.Printf("%s", time.Since(start))
	time.Sleep(time.Second * 2)
	fmt.Println("\nEnd")

	// c := make(chan bool, 2)

	// for true {
	// 	go func() {
	// 		fmt.Println(time.Now().Second())
	// 		time.Sleep(2 * time.Second)
	// 		<-c
	// 	}()

	// 	c <- true
	// }

}
