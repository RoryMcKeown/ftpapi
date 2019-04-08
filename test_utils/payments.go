package test_utils

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/rorymckeown/ftpapi/db"
	"github.com/rorymckeown/ftpapi/swagger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gotest.tools/assert"
)

func ClonePaymentWithNewId(payment *swagger.Payment) *swagger.Payment {
	return ClonePaymentAndSetId(payment, uuid.UUID.String(uuid.New()))
}

func ClonePaymentAndSetId(payment *swagger.Payment, id string) *swagger.Payment {
	clone := swagger.Payment{}
	copier.Copy(&clone, payment)
	clone.Id = id
	return &clone
}

func ClonePayment(payment *swagger.Payment) *swagger.Payment {
	clone := swagger.Payment{}
	copier.Copy(&clone, payment)
	return &clone
}

func GetSamplePayment(t *testing.T) *swagger.Payment {
	//This is kind of filthy.
	//When the tests run they use their local 'testdata', rather than a 'test_utils/testdata' dir.
	//So putting this file relative to the root project
	path := filepath.Join("../testdata", "payment.json")
	bytes, errRead := ioutil.ReadFile(path)
	assert.NilError(t, errRead)

	p, errDecode := db.Decode(bytes)
	assert.NilError(t, errDecode)
	return p
}

func AssertPaymentInPaymentsList(t *testing.T, payment *swagger.Payment, list []*swagger.Payment) {
	for i := range list {
		if payment.Id == list[i].Id {
			assert.DeepEqual(t, payment, list[i])
			return
		}
	}

	assert.Assert(t, "payment with id was not found in list %s", payment.Id)
}
