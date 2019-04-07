package test_utils_test

import (
	test_utils "github.com/rorymckeown/ftpapi/test_utils"

	"testing"

	"github.com/rorymckeown/ftpapi/swagger"

	"gotest.tools/assert"
)

func TestClonePayment(t *testing.T) {

	p1 := test_utils.GetSamplePayment(t)

	assert.Equal(t, p1.Id, "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43")
	assert.Equal(t, p1.OrganisationId, "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb")

	p2 := test_utils.ClonePayment(p1)

	assert.Assert(t, p1.OrganisationId == p2.OrganisationId)

	assert.Equal(t, p2.Id, "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43")
	assert.Equal(t, p2.OrganisationId, "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb")

	p3 := test_utils.ClonePaymentWithNewId(p1)
	assert.Assert(t, p1.OrganisationId == p3.OrganisationId)
	assert.Assert(t, p1.Id != p3.Id)
}

func TestAssertPaymentInList(t *testing.T) {
	basePayment := test_utils.GetSamplePayment(t)

	numPayments := 10
	payments := make([]*swagger.Payment, numPayments)

	for i := range payments {
		payments[i] = test_utils.ClonePaymentWithNewId(basePayment)
	}

	for i := range payments {
		test_utils.AssertPaymentInPaymentsList(t, payments[i], payments)
	}
}
