package db

import (
	swagger "github.com/rorymckeown/ftpapi/swagger"

	"encoding/json"

	leveldb "github.com/syndtr/goleveldb/leveldb"
	opt "github.com/syndtr/goleveldb/leveldb/opt"
)

func Encode(payment *swagger.Payment) ([]byte, error) {
	return json.Marshal(payment)
}

func Decode(data []byte) (*swagger.Payment, error) {
	payment := swagger.Payment{}
	err := json.Unmarshal(data, &payment)
	return &payment, err
}

func OpenDB(path string) (*leveldb.DB, error) {

	return leveldb.OpenFile(path, nil)
}

func CloseDB(db *leveldb.DB) error {
	return db.Close()
}

func PutPayment(db *leveldb.DB, payment *swagger.Payment) error {
	data, err := Encode(payment)

	if err != nil {
		return err
	}

	opts := opt.WriteOptions{Sync: true}

	return db.Put([]byte(payment.Id), data, &opts)
}

func GetPayment(db *leveldb.DB, id string) (*swagger.Payment, error) {

	data, err := db.Get([]byte(id), nil)

	if err != nil {
		return nil, err
	}

	return Decode(data)
}

func DeletePayment(db *leveldb.DB, id string) error {
	return db.Delete([]byte(id), nil)
}

func ListPayments(db *leveldb.DB) []*swagger.Payment {

	result := []*swagger.Payment{}

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		value := iter.Value()

		payment, err := Decode(value)

		if err == nil {
			result = append(result, payment)
		} else {
			//TODO: Log something?
		}
	}

	return result
}
