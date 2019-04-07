package db_test

import (
	"github.com/rorymckeown/ftpapi/db"
	"github.com/rorymckeown/ftpapi/swagger"
	test_utils "github.com/rorymckeown/ftpapi/test_utils"

	"os"

	"testing"

	"gotest.tools/assert"
)

func TestEncodeDecodePayment(t *testing.T) {

	payment := getSamplePayment(t)

	bytes2, errEncode := db.Encode(payment)
	assert.NilError(t, errEncode)

	payment2, errDecode := db.Decode(bytes2)
	assert.NilError(t, errDecode)
	assert.DeepEqual(t, payment, payment2)
}

func TestPostGetDeletePayment(t *testing.T) {

	path := "/tmp/TestPostGetDeletePayment/"
	levelDb, errDbOpen := db.OpenDB(path)
	assert.NilError(t, errDbOpen)
	defer db.CloseDB(levelDb)
	defer os.RemoveAll(path)

	payment := getSamplePayment(t)

	putError := db.PutPayment(levelDb, payment)
	assert.NilError(t, putError)

	payment2, errGet := db.GetPayment(levelDb, payment.Id)

	assert.Assert(t, payment2 != nil)
	assert.NilError(t, errGet)
	assert.DeepEqual(t, payment, payment2)

	errDelete := db.DeletePayment(levelDb, payment.Id)
	assert.NilError(t, errDelete)

	payment3, errGet3 := db.GetPayment(levelDb, payment.Id)
	assert.Assert(t, payment3 == nil)
	assert.Assert(t, errGet3 != nil)
}

func TestPutManyAndListPayments(t *testing.T) {
	path := "/tmp/TestPutManyAndListPayments/"
	levelDb, errDbOpen := db.OpenDB(path)
	assert.NilError(t, errDbOpen)
	defer db.CloseDB(levelDb)
	defer os.RemoveAll(path)

	numPayments := 10
	basePayment := getSamplePayment(t)

	payments := make([]*swagger.Payment, numPayments)

	for i := range payments {
		payments[i] = test_utils.ClonePaymentWithNewId(basePayment)
		putError := db.PutPayment(levelDb, payments[i])
		assert.NilError(t, putError)
	}

	paymentsList := db.ListPayments(levelDb)

	assert.Equal(t, numPayments, len(paymentsList))

	for i := range payments {
		test_utils.AssertPaymentInPaymentsList(t, payments[i], paymentsList)
	}

}

func getSamplePayment(t *testing.T) *swagger.Payment {
	return test_utils.GetSamplePayment(t)
}
