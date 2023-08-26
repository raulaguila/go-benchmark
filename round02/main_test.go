package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const paramNG string = "paramNG"

type bench struct {
	name           string
	param          string
	resultExpected int
	url            string
	body           []byte
}

func makeRequest(param, path string, body []byte) int {
	res, err := http.Post(fmt.Sprintf("%s/%s", path, param), "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	return res.StatusCode
}

func BenchmarkWebFramework(b *testing.B) {
	const urlAtreugo string = "http://localhost:8086/bench"
	const urlGoFiber string = "http://localhost:8083/bench"
	const urlGoHttp string = "http://localhost:8085/bench"

	const urlGoGin string = "http://localhost:8080/bench"
	const urlGorillaMux string = "http://localhost:8081/bench"
	const urlGoChi string = "http://localhost:8082/bench"
	const urlEcho string = "http://localhost:8084/bench"

	obj := ObjectExample{"testId", "testName", 53}
	jsn, _ := json.Marshal(obj)

	var tests []bench = []bench{
		{name: "tests_right_param_ATREUGO", param: paramOK, resultExpected: http.StatusOK, url: urlAtreugo, body: jsn},
		{name: "tests_wrong_param_ATREUGO", param: paramNG, resultExpected: http.StatusBadRequest, url: urlAtreugo, body: jsn},
		{name: "tests_wrong_body_ATREUGO", param: paramOK, resultExpected: http.StatusBadRequest, url: urlAtreugo, body: nil},

		{name: "tests_right_param_GOFIBER", param: paramOK, resultExpected: http.StatusOK, url: urlGoFiber, body: jsn},
		{name: "tests_wrong_param_GOFIBER", param: paramNG, resultExpected: http.StatusBadRequest, url: urlGoFiber, body: jsn},
		{name: "tests_wrong_body_GOFIBER", param: paramOK, resultExpected: http.StatusBadRequest, url: urlGoFiber, body: nil},

		{name: "tests_right_param_HTTPSERVERMUX", param: paramOK, resultExpected: http.StatusOK, url: urlGoHttp, body: jsn},
		{name: "tests_wrong_param_HTTPSERVERMUX", param: paramNG, resultExpected: http.StatusBadRequest, url: urlGoHttp, body: jsn},
		{name: "tests_wrong_body_HTTPSERVERMUX", param: paramOK, resultExpected: http.StatusBadRequest, url: urlGoHttp, body: nil},

		{name: "tests_right_param_GINGONIC", param: paramOK, resultExpected: http.StatusOK, url: urlGoGin, body: jsn},
		{name: "tests_wrong_param_GINGONIC", param: paramNG, resultExpected: http.StatusBadRequest, url: urlGoGin, body: jsn},
		{name: "tests_wrong_body_GINGONIC", param: paramOK, resultExpected: http.StatusBadRequest, url: urlGoGin, body: nil},

		{name: "tests_right_param_GOCHI", param: paramOK, resultExpected: http.StatusOK, url: urlGoChi, body: jsn},
		{name: "tests_wrong_param_GOCHI", param: paramNG, resultExpected: http.StatusBadRequest, url: urlGoChi, body: jsn},
		{name: "tests_wrong_body_GOCHI", param: paramOK, resultExpected: http.StatusBadRequest, url: urlGoChi, body: nil},

		{name: "tests_right_param_GORILLAMUX", param: paramOK, resultExpected: http.StatusOK, url: urlGorillaMux, body: jsn},
		{name: "tests_wrong_param_GORILLAMUX", param: paramNG, resultExpected: http.StatusBadRequest, url: urlGorillaMux, body: jsn},
		{name: "tests_wrong_body_GORILLAMUX", param: paramOK, resultExpected: http.StatusBadRequest, url: urlGorillaMux, body: nil},

		{name: "tests_right_param_ECHO", param: paramOK, resultExpected: http.StatusOK, url: urlEcho, body: jsn},
		{name: "tests_wrong_param_ECHO", param: paramNG, resultExpected: http.StatusBadRequest, url: urlEcho, body: jsn},
		{name: "tests_wrong_body_ECHO", param: paramOK, resultExpected: http.StatusBadRequest, url: urlEcho, body: nil},
	}

	startFrameworks()
	time.Sleep(5 * time.Second)

	for u, test := range tests {
		b.Run(test.name, func(bf *testing.B) {
			bf.ReportAllocs()
			for i := 0; i < bf.N; i++ {
				if result := makeRequest(test.param, test.url, test.body); result != test.resultExpected {
					bf.Errorf("%v - Result not expected, expected=%v, received=%v", test.name, test.resultExpected, result)
					return
				}
			}
		})
		if (u+1)%3 == 0 {
			println("")
		}
	}
}
