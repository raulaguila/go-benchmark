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

func makeRequest(param, path string, body []byte) int {
	response, err := http.Post(fmt.Sprintf("%s/%s", path, param), "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	return response.StatusCode
}

func BenchmarkWebFramework(b *testing.B) {
	obj := &ObjectExample{"testId", "testName", 53}
	objJson, _ := json.Marshal(obj)
	objEmpty := &ObjectExample{}

	tests := []bench{}
	for _, test := range [][2]string{
		{"ATREUGO", atreugoPort},
		{"HTTPMUX", httpPort},
		{"GINGONIC", ginPort},
		{"GOFIBER", fiberPort},
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
				if responseStatus := makeRequest(test.requestParam, test.requestURL, test.requestBody); responseStatus != test.expectedStatus {
					bf.Errorf("%v - Result not expected: \n\t\tstatusExpected=%v, statusReceived=%v \n\n", test.testName, test.expectedStatus, responseStatus)
					return
				}
			}
		})
		if (u+1)%3 == 0 {
			println("")
		}
	}
}
