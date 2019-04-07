package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	db "github.com/rorymckeown/ftpapi/db"
	swagger "github.com/rorymckeown/ftpapi/swagger"

	leveldb "github.com/syndtr/goleveldb/leveldb"
)

func GetPaymentRoutes(db *leveldb.DB) []Route {
	return []Route{

		Route{
			"PaymentsPaymentIdPut",
			strings.ToUpper("Put"),
			"/v2/payments/{payment_id}",
			func(w http.ResponseWriter, r *http.Request) {
				PaymentsPaymentIdPut(db, w, r)
			},
		},

		Route{
			"PaymentsGet",
			strings.ToUpper("Get"),
			"/v2/payments",
			func(w http.ResponseWriter, r *http.Request) {
				PaymentsGet(db, w, r)
			},
		},

		Route{
			"PaymentsPaymentIdDelete",
			strings.ToUpper("Delete"),
			"/v2/payments/{payment_id}",
			func(w http.ResponseWriter, r *http.Request) {
				PaymentsPaymentIdDelete(db, w, r)
			},
		},

		Route{
			"PaymentsPaymentIdGet",
			strings.ToUpper("Get"),
			"/v2/payments/{payment_id}",
			func(w http.ResponseWriter, r *http.Request) {
				PaymentsPaymentIdGet(db, w, r)
			},
		},

		Route{
			"PaymentsPost",
			strings.ToUpper("Post"),
			"/v2/payments",
			func(w http.ResponseWriter, r *http.Request) {
				PaymentsPost(db, w, r)
			},
		},
	}
}

func PaymentsGet(ldb *leveldb.DB, w http.ResponseWriter, r *http.Request) {

	paymentsList := db.ListPayments(ldb)

	data, err := json.Marshal(swagger.PaymentList{Data: paymentsList})

	if err != nil {
		//TODO: Add an error struct
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

func PaymentsPaymentIdDelete(ldb *leveldb.DB, w http.ResponseWriter, r *http.Request) {

	paymentId := paymentIdFromRequest(r)

	_, errGetPayment := db.GetPayment(ldb, paymentId)

	if errGetPayment != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	errDelete := db.DeletePayment(ldb, paymentId)

	if errDelete != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func paymentIdFromRequest(r *http.Request) string {
	parts := strings.Split(r.URL.Path, "/")
	return parts[len(parts)-1]
}

func PaymentsPaymentIdGet(ldb *leveldb.DB, w http.ResponseWriter, r *http.Request) {

	paymentId := paymentIdFromRequest(r)

	payment, errGetPayment := db.GetPayment(ldb, paymentId)

	if errGetPayment != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, errEncode := db.Encode(payment)

	if errEncode != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}

func PaymentsPaymentIdPut(ldb *leveldb.DB, w http.ResponseWriter, r *http.Request) {

	paymentId := paymentIdFromRequest(r)

	newPayment, err := paymentFromRequest(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if paymentId != newPayment.Id {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, errGetPayment := db.GetPayment(ldb, newPayment.Id)

	if errGetPayment != nil {
		//TODO: Make a decent error struct ?
		w.WriteHeader(http.StatusNotFound)
		return
	}

	errPut := db.PutPayment(ldb, newPayment)

	if errPut != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func paymentFromRequest(r *http.Request) (*swagger.Payment, error) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	return db.Decode(body)
}

func PaymentsPost(ldb *leveldb.DB, w http.ResponseWriter, r *http.Request) {

	newPayment, err := paymentFromRequest(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	previousPayment, errGetPayment := db.GetPayment(ldb, newPayment.Id)

	if errGetPayment == nil && previousPayment != nil {
		//TODO: Make a decent error struct ?
		w.WriteHeader(http.StatusConflict)
		return
	}

	putErr := db.PutPayment(ldb, newPayment)

	if putErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
