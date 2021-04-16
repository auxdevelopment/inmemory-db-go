package main

import (
	"flag"
	"fmt"
	"net/http"

	"auxdev.at/storage-node/server"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	server := server.NewServer()

	portFlag := flag.Int("port", 8080, "Port on which the storage server will listen")
	flag.Parse()

	router.HandleFunc("/values", server.StorePairHandler).Methods("PUT")
	router.HandleFunc("/values/{key}", server.GetPairHandler).Methods("GET")
	router.HandleFunc("/values", server.GetAllPairsHandler).Methods("GET")
	router.HandleFunc("/status", server.GetStatusHandler).Methods("GET")

	fmt.Println("Starting server")
	err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), router)
	if err != nil {
		fmt.Println("Could not start server.")
	}
}
