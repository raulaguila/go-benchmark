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
	testName       string
	requestParam   string
	requestURL     string
	requestBody    []byte
	expectedStatus int
	expectedBody   *ObjectExample
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
			bench{testName: "all_good_" + test[0], requestParam: paramOK, expectedStatus: http.StatusOK, requestURL: url, requestBody: objJson, expectedBody: obj},
			bench{testName: "wrong_param_" + test[0], requestParam: paramNG, expectedStatus: http.StatusBadRequest, requestURL: url, requestBody: objJson, expectedBody: objEmpty},
			bench{testName: "wrong_body_" + test[0], requestParam: paramOK, expectedStatus: http.StatusBadRequest, requestURL: url, requestBody: nil, expectedBody: objEmpty},
		)
	}

	startFrameworks()
	time.Sleep(2 * time.Second)

	for u, test := range tests {
		b.Run(test.testName, func(bf *testing.B) {
			bf.ReportAllocs()
			for i := 0; i < bf.N; i++ {
				if responseStatus, responseBody := makeRequest(test.requestParam, test.requestURL, test.requestBody); responseStatus != test.expectedStatus || *responseBody != *test.expectedBody {
					bf.Errorf("%v - Result not expected: \n\t\tstatusExpected=%v, statusReceived=%v, \n\t\tbodyExpected=%v, bodyReceived=%v \n\n", test.testName, test.expectedStatus, responseStatus, *test.expectedBody, *responseBody)
					return
				}
			}
		})
		if (u+1)%3 == 0 {
			println("")
		}
	}
}
