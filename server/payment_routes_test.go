package server_test

import (
	"bytes"

	"encoding/json"
	"io/ioutil"
	http "net/http"
	"os"
	"strings"
	"testing"

	"github.com/rorymckeown/ftpapi/server"
	swagger "github.com/rorymckeown/ftpapi/swagger"
	test_utils "github.com/rorymckeown/ftpapi/test_utils"

	"gotest.tools/assert"
)

// The big flaw in these tests is that http servers are started and not stopped.
// This is due to some implementation of the http server and the interaction with the test framework I really know nothing about.
// Attempting to reuse the same port causes 'binding' issues after the first test
// So each test gets a new server and a new db
// This means I dont get data collisions between tests

func StartServer(endpoint string, port int) {
	exitChan := make(chan int)
	server.StartServer(endpoint, port, exitChan)
}
func TestEmptyPaymentListGet(t *testing.T) {
	StartServer("/tmp/TestEmptyPaymentListGet", 8080)
	defer os.RemoveAll("/tmp/TestPaymentListGet")

	paymentList := getPaymentList(t, "http://localhost:8080/v2/payments")

	assert.Assert(t, 0 == len(paymentList))
}

func TestPostGetDeletePayment(t *testing.T) {

	StartServer("/tmp/TestPostPayment", 8081)
	defer os.RemoveAll("/tmp/TestPostPayment")

	payment := test_utils.GetSamplePayment(t)

	postPayment(t, "http://localhost:8081/v2/payments", payment)

	paymentList := getPaymentList(t, "http://localhost:8081/v2/payments")

	test_utils.AssertPaymentInPaymentsList(t, payment, paymentList)

	getPayment := getPayment(t, "http://localhost:8081/v2/payments", payment.Id)

	assert.DeepEqual(t, getPayment, payment)

	//very odd that this is not part of the standard api
	deleteReq, errDeleteReq := http.NewRequest("DELETE", joinUrl("http://localhost:8081/v2/payments", payment.Id), nil)
	assert.NilError(t, errDeleteReq)
	deleteResp, errDeletePayment := http.DefaultClient.Do(deleteReq)
	assert.NilError(t, errDeletePayment)
	assert.Equal(t, deleteResp.StatusCode, 200)

	emptyPaymentList := getPaymentList(t, "http://localhost:8081/v2/payments")

	assert.Equal(t, len(emptyPaymentList), 0)
}

func TestPutPayment(t *testing.T) {
	StartServer("/tmp/TestPutPayment", 8082)
	defer os.RemoveAll("/tmp/TestPutPayment")
	endpoint := "http://localhost:8082/v2/payments"

	payment := test_utils.GetSamplePayment(t)

	postPayment(t, endpoint, payment)

	updatePayment := test_utils.ClonePayment(payment)
	updatePayment.OrganisationId = "UpdatedOrgId"

	assert.Assert(t, updatePayment.OrganisationId != payment.OrganisationId)

	updateBody, errMarshal := json.Marshal(updatePayment)
	assert.NilError(t, errMarshal)

	putResp, putErr := putReq(joinUrl(endpoint, payment.Id), updateBody)

	assert.Equal(t, putResp.StatusCode, 200)
	bodyBytes := getBody(putResp)
	assert.Equal(t, 0, len(bodyBytes))
	assert.NilError(t, putErr)

	getPayment := getPayment(t, endpoint, payment.Id)
	assert.DeepEqual(t, updatePayment, getPayment)
	assert.Assert(t, getPayment.OrganisationId != payment.OrganisationId)
}

func TestPostGetMany(t *testing.T) {
	StartServer("/tmp/TestPostGetMany", 8083)
	defer os.RemoveAll("/tmp/TestPostGetMany")
	endpoint := "http://localhost:8083/v2/payments"

	numPayments := 10
	payments := make([]*swagger.Payment, numPayments)
	basePayment := test_utils.GetSamplePayment(t)

	for i := range payments {
		payments[i] = test_utils.ClonePaymentWithNewId(basePayment)
		postPayment(t, endpoint, payments[i])
	}

	paymentsList := getPaymentList(t, endpoint)

	assert.Equal(t, numPayments, len(paymentsList))

	for i := range payments {
		test_utils.AssertPaymentInPaymentsList(t, payments[i], paymentsList)
	}
}

func TestPostErrors(t *testing.T) {
	StartServer("/tmp/TestPostErrors", 8084)
	defer os.RemoveAll("/tmp/TestPostErrors")
	endpoint := "http://localhost:8084/v2/payments"

	//Post bad json
	body := []byte("badjson")
	resp, postErr := http.Post(endpoint, "application/json", bytes.NewBuffer(body))

	assert.Equal(t, resp.StatusCode, 400)
	assert.NilError(t, postErr)

	//Post same payment twice
	basePayment := test_utils.GetSamplePayment(t)
	postPayment(t, endpoint, basePayment)

	body, errMarshal := json.Marshal(basePayment)
	assert.NilError(t, errMarshal)

	resp, secondPostErr := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	assert.Equal(t, resp.StatusCode, 409)
	assert.NilError(t, secondPostErr)
}

func TestGetErrors(t *testing.T) {
	StartServer("/tmp/TestGetErrors", 8085)
	defer os.RemoveAll("/tmp/TestGetErrors")
	endpoint := "http://localhost:8085/v2/payments"

	//Get non existing payment
	getPaymentResp, errGetPayment := http.Get(joinUrl(endpoint, "nonexistingpayment"))
	assert.NilError(t, errGetPayment)
	assert.Equal(t, getPaymentResp.StatusCode, 404)

	//Get payment with wrong url
	payment := test_utils.GetSamplePayment(t)
	postPayment(t, endpoint, payment)
	getPayment(t, endpoint, payment.Id)

	getPaymentResp2, errGetPayment2 := http.Get(joinUrl(endpoint, joinUrl(payment.Id, "foo")))
	assert.NilError(t, errGetPayment2)
	assert.Equal(t, getPaymentResp2.StatusCode, 404)
}

func TestPutErrors(t *testing.T) {
	StartServer("/tmp/TestPostErrors", 8086)
	defer os.RemoveAll("/tmp/TestPostErrors")
	endpoint := "http://localhost:8086/v2/payments"

	//Put non existing payment
	payment := test_utils.GetSamplePayment(t)
	paymentBody, errMarshal := json.Marshal(payment)
	assert.NilError(t, errMarshal)

	getPaymentResp, errGetPayment := http.Get(joinUrl(endpoint, payment.Id))
	assert.NilError(t, errGetPayment)
	assert.Equal(t, getPaymentResp.StatusCode, 404)

	putResp, putErr := putReq(joinUrl(endpoint, payment.Id), paymentBody)
	assert.NilError(t, putErr)
	assert.Equal(t, putResp.StatusCode, 404)

	//Make a target to update incorrectly
	postPayment(t, endpoint, payment)
	getPaymentResp2, errGetPayment2 := http.Get(joinUrl(endpoint, payment.Id))
	assert.NilError(t, errGetPayment2)
	assert.Equal(t, getPaymentResp2.StatusCode, 200)

	//Put invalid body
	putInvalidBodyResp, errPutInvalidBodyResp := putReq(joinUrl(endpoint, payment.Id), []byte("invalidjson"))
	assert.NilError(t, errPutInvalidBodyResp)
	assert.Equal(t, putInvalidBodyResp.StatusCode, 400)

	//Put invalid id
	payment2 := test_utils.ClonePaymentWithNewId(payment)
	payment2Body, errMarshal := json.Marshal(payment2)
	assert.NilError(t, errMarshal)
	putInvalidIdResp, errPutInvalidId := putReq(joinUrl(endpoint, payment.Id), payment2Body)
	assert.NilError(t, errPutInvalidId)
	assert.Equal(t, putInvalidIdResp.StatusCode, 400)

}

func joinUrl(prefix string, end string) string {
	return strings.Join([]string{prefix, end}, "/")
}

func unmarshalPaymentList(bytes []byte) *swagger.PaymentList {
	paymentList := swagger.PaymentList{}
	json.Unmarshal(bytes, &paymentList)
	return &paymentList
}

func getBody(resp *http.Response) []byte {
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func getPaymentList(t *testing.T, url string) []*swagger.Payment {

	resp, err := http.Get(url)
	assert.Equal(t, resp.StatusCode, 200)
	assert.NilError(t, err)

	bodyBytes := getBody(resp)

	return unmarshalPaymentList(bodyBytes).Data
}

func postPayment(t *testing.T, url string, payment *swagger.Payment) {
	body, errMarshal := json.Marshal(payment)
	assert.NilError(t, errMarshal)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	assert.Equal(t, resp.StatusCode, 200)
	bodyBytes := getBody(resp)
	assert.Equal(t, 0, len(bodyBytes))
	assert.NilError(t, err)
}

func getPayment(t *testing.T, url string, id string) *swagger.Payment {
	getPaymentResp, errGetPayment := http.Get(joinUrl(url, id))

	assert.NilError(t, errGetPayment)

	getPaymentBodyBytes := getBody(getPaymentResp)
	getPayment := swagger.Payment{}
	errUnmarshal := json.Unmarshal(getPaymentBodyBytes, &getPayment)
	assert.NilError(t, errUnmarshal)
	return &getPayment
}

func putReq(endpoint string, body []byte) (*http.Response, error) {
	putReq, _ := http.NewRequest("PUT", endpoint, bytes.NewBuffer(body))
	return http.DefaultClient.Do(putReq)
}
