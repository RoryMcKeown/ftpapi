package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	db "github.com/rorymckeown/ftpapi/db"
)

func StartServer(dbPath string, port int, exitChan chan<- int) *http.Server {

	log.Printf("Server started")

	leveldb, errOpenDB := db.OpenDB(dbPath)

	if errOpenDB != nil {
		panic(errOpenDB)
	}

	portStr := strconv.Itoa(port)

	router := CreateRouterForRoutes(GetPaymentRoutes(leveldb))

	srv := &http.Server{
		Addr:    strings.Join([]string{":", portStr}, ""),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			db.CloseDB(leveldb)
			log.Fatal(err)
			exitChan <- 1
		}
	}()

	return srv
}
