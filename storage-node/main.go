package main

import (
	"flag"
	"fmt"
	"net/http"

	"auxdev.at/storage-node/server"
	"github.com/gorilla/mux"
)

const (
	defaultPort       = 8080
	defaultConfigPath = ""
)

func main() {
	router := mux.NewRouter()
	server := server.NewServer()

	configFlag := flag.String("conf", "", "Path to configuration file.")
	portFlag := flag.Int("port", 8080, "Port on which the storage server will listen")
	flag.Parse()

	if *configFlag != defaultConfigPath {
		// TODO
	}

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
