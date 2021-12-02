package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/3dsinteractive/wrkgo"
)

func main() {

	config := &wrkgo.LoadTestConfig{
		BaseURL:         "http://localhost:8088",
		ConcurrentUsers: 50,
		RunDuration:     time.Second * 10,
		DebugError:      true,
		DebugRequest:    false,
		DebugResponse:   false,
	}

	templates := []*wrkgo.LoadTestTemplate{
		{
			ID:      "0",
			URLPath: "/profile",
			Timeout: time.Second * 6,
			Method:  "GET",
			Headers: map[string]string{
				"Content-Type":  "application/json; charset=UTF-8",
				"Authorization": "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyODMxNTYsInVzZXJuYW1lIjoiZGV2MDEiLCJuYW1lIjoiZGV2IGRldiJ9.Ns_Y5dRP503_hQaPYarORyzK-hinp0Snf8KpLipoikprKNmDwj2WaroUaligd1qt7OENh6RYOMA6WiHds0PGyS0URBb4tHxkbYJhN23SPaR8D8-tJb7h7EpNaEUFSkLific_yn9WCMzMTEkeKd-jZRKCtOgiMdB4AhJZ2SpsRu5Z_UaN3ycrRscF-FIMtVOMQU01Zg18Q9NNAOZi79Ll0skdITt0UZY9xWjdrBKlHy_WMYhcvlYNMhZ9UzlB8BoJDNbvx1MckJCGP2iF128cU25EmqFYsGnIZE22UCEF94eObPyFF6QoUhfMbTbjRB9_0fConNfPT6jG3JLY9A6vtpFjJnhnZ5nxFx3W0PkRrEy_1hllcFwyVVQGZSD_9yUyAlJTyvksMPHZh3vm4qxuWmZB-A9vF__m3wTijuZHPaAHqYKsr-WNLhLcacfsT5wyEGWxHAbaQPGwtWz4zs_HpFDYHFBnr2OdsukCOJCP1yb8QGAehvY_fh2OgnYO_IADxsnj-jzy4Ng6H9-yzcHuxTb4pxkoRDu28l-uTLV09Xdhydsa3sZcyFq11GjYYpMHybEpHuFmHtOVCBljV8yUUCRBv7Ze_iLDiRJWH6KNFZyPIJz-Hg7PmelCbvIbltanwKyxizrlKEByd-v6o2cq1oJ939fz4qmSlKhdFbFsGt8",
			},
		},
	}

	reqSetupHandler := func(tmpl *wrkgo.LoadTestTemplate, req *wrkgo.LoadTestRequest, prevResp *wrkgo.LoadTestResponse) error {
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
