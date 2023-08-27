package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	localhost string = "http://localhost"
	paramNG   string = "paramNG"
)

type bench struct {
	name           string
	param          string
	url            string
	body           []byte
	statusExpected int
	bodyExpected   *ObjectExample
}

func makeRequest(param, path string, body []byte) (int, *ObjectExample) {
	response, err := http.Post(fmt.Sprintf("%s/%s", path, param), "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	objReceived := new(ObjectExample)
	json.NewDecoder(response.Body).Decode(objReceived)

	return response.StatusCode, objReceived
}

func BenchmarkWebFramework(b *testing.B) {
	obj := &ObjectExample{"testId", "testName", 53}
	objJson, _ := json.Marshal(obj)
	objEmpty := &ObjectExample{}

	tests := []bench{}
	for _, test := range [][2]string{
		{"ATREUGO", atreugoPort},
		{"GOFIBER", fiberPort},
		{"HTTPMUX", httpPort},
		{"GINGONIC", ginPort},
		{"GOCHI", chiPort},
		{"GORILLAMUX", gorillaPort},
		{"GOECHO", echoPort},
	} {
		url := localhost + test[1] + endpoint
		tests = append(tests,
			bench{name: "all_good_" + test[0], param: paramOK, statusExpected: http.StatusOK, url: url, body: objJson, bodyExpected: obj},
			bench{name: "wrong_param_" + test[0], param: paramNG, statusExpected: http.StatusBadRequest, url: url, body: objJson, bodyExpected: objEmpty},
			bench{name: "wrong_body_" + test[0], param: paramOK, statusExpected: http.StatusBadRequest, url: url, body: nil, bodyExpected: objEmpty},
		)
	}

	startFrameworks()
	time.Sleep(1 * time.Second)

	for u, test := range tests {
		b.Run(test.name, func(bf *testing.B) {
			bf.ReportAllocs()
			for i := 0; i < bf.N; i++ {
				if responseStatus, responseBody := makeRequest(test.param, test.url, test.body); responseStatus != test.statusExpected || *responseBody != *test.bodyExpected {
					bf.Errorf("%v - Result not expected: \n\t\tstatusExpected=%v, statusReceived=%v, \n\t\tbodyExpected=%v, bodyReceived=%v \n\n", test.name, test.statusExpected, responseStatus, *test.bodyExpected, *responseBody)
					return
				}
			}
		})
		if (u+1)%3 == 0 {
			println("")
		}
	}
}
