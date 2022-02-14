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
		ConcurrentUsers: 1,
		RunDuration:     time.Second * 1,
		DebugError:      true,
		DebugRequest:    false,
		DebugResponse:   false,
	}

	templates := []*wrkgo.LoadTestTemplate{
		{
			ID:      "0",
			URLPath: "/merchant/23xK48ZSaDPzoxZVXIbV8w6kFVw/inventory",
			Timeout: time.Second * 6,
			Method:  "POST",
			Headers: map[string]string{
				"Content-Type":  "application/json; charset=UTF-8",
				"Authorization": "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDUzNDQ5NDksInVzZXJuYW1lIjoiZGV2MDEiLCJuYW1lIjoiZGV2IGRldiJ9.qTWCinNiknDUDafqSMr3HM8qjxcnlV5r8ckI4QD8BTkhB7XsVpVgwLiWBUcGMH26TBzf7TD9puu4Jqt0SByrM34RY98NmVTw3qQfJ19EkmqGcH8o7fhxULqs2GW8pgDqANsYEW1n_ibZ5R0sUjmxixpPUgOKZMSl2Wdtw5pipZyQmzznOR14Nky7szSwFoT8gXB0ps3EALYTKZYbkyRJ-4SFLZppwWq3QLJNaIePmR2N7wIAQlzWdLLSen0onk248DisDB5taM5aCHR7aBFNu02bv2U3b8GyRTJ2pWsA4iWPlOM3QeufuU_MYgw-ThQSnWl--KXJhHxaC8bngpUnqX9h1qTEie7sx2HdaTjLuuKyqegyCgia5yCyxRkf7-NGxXH5vuP2dI7FtV-JHLgSse1isjTHXfM9ZYwAleuBtwbSGwD126pIr6-KL7tk1I0cKtGgA4dm0_MQFPmLVSFXV6aAebuiYpPRFBl950mU5qe4s-Q7Ox_1lANPKZfZz0OaMTaRyePtX39IrQIBjPr0T6z_CVIS8uT_mNoO14ojwpU4azM0Y0VrlD5epp0MplZpDHq34dd9iPy9q5pCiSNmfp3kOjmg3U3dBTqE4EcT6QoSmN79vB-Q5ZtDFs2LC5-WUhVPdNwk_IcKAyWKavNkcuQCRnjm1gKebfD_tuwZ_8E",
			},
		},
	}

	runNumber := 1
	reqSetupHandler := func(tmpl *wrkgo.LoadTestTemplate, req *wrkgo.LoadTestRequest, prevResp *wrkgo.LoadTestResponse) error {

		runNumberTxt := strconv.Itoa(runNumber)
		input := map[string]interface{}{
			"itemSku":      "dev" + runNumberTxt,
			"merchantId":   "23twO9nFtgsLGAuQ9JXPzi3C65N",
			"categoryGuid": "cat1234",
			"lineNumber":   1,
			"price":        99.0,
			"recommended":  false,
			"haveImage":    false,
			"activated":    true,
			"name1":        "Name " + runNumberTxt,
			"description1": "Description " + runNumberTxt,
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
