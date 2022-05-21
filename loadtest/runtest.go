package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/3dsinteractive/wrkgo"
)

func main() {

	config := &wrkgo.LoadTestConfig{
		BaseURL:         "http://localhost:8088",
		ConcurrentUsers: 700,
		RunDuration:     time.Second * 60 * 10,
		DebugError:      true,
		DebugRequest:    false,
		DebugResponse:   false,
	}

	templates := []*wrkgo.LoadTestTemplate{
		{
			ID:      0,
			URLPath: "/transaction",
			Timeout: time.Second * 6,
			Method:  "POST",
			Headers: map[string]string{
				"Content-Type":  "application/json; charset=UTF-8",
				"Authorization": "Bearer 26j5fKC8fQ0dFNza2pq9degmhTS",
			},
		},
	}

	runNumber := 1
	reqSetupHandler := func(tmpl *wrkgo.LoadTestTemplate, req *wrkgo.LoadTestRequest, prevResp *wrkgo.LoadTestResponse) error {
		runNumberTxt := strconv.Itoa(runNumber)
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{
					"inventoryId":  "inv",
					"itemSku":      "sku" + runNumberTxt,
					"categoryGuid": fmt.Sprintf("cat-%d", rand.Int()),
					"lineNumber":   1,
					"price":        rand.Float64(),
					"qty":          rand.Int(),
				},
			},
		}
		req.SetBodyJSON(input)
		runNumber++
		return nil
	}

	lt := wrkgo.NewLoadTest()
	err := lt.Run(config, templates, reqSetupHandler)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func RandomMinMax(min int, max int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(max-min+1) + min
}
