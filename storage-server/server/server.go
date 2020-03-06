package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"auxdev.at/storage-server/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	hashMap *storage.HashMap
	logger  *log.Logger
}

func NewServer() *Server {
	return &Server{
		storage.NewHashMap(),
		log.New(
			os.Stdout,
			"Storage: ",
			log.LUTC|log.Ltime|log.Ldate|log.Lmicroseconds,
		),
	}
}

func (server *Server) StorePairHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var pair storage.KeyValuePair

	if err := decoder.Decode(&pair); err != nil {
		server.logger.Println(err)

		w.WriteHeader(500)
		w.Write([]byte("{ \"message\": \"Internal server error while processing request body\" }\n"))

		return
	}

	if err := server.hashMap.Put(pair); err != nil {
		server.logger.Println(err)

		w.WriteHeader(500)
		w.Write([]byte("{ \"message\": \"Internal server error while trying to store pair\" }\n"))

		return
	}

	server.logger.Printf("Stored %s", pair)

	defer r.Body.Close()

	w.WriteHeader(200)
	w.Write([]byte("{ \"message\": \"OK\" }\n"))
}

func (server *Server) GetPairHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	params := mux.Vars(r)
	key := params["key"]

	result, err := server.hashMap.Get(key)

	if err != nil {
		server.logger.Println(err)

		w.WriteHeader(404)
		w.Write([]byte("{ \"message\": \"Requested resource could not be found\" }\n"))

		return
	}

	responseBody, _ := json.Marshal(result)

	w.WriteHeader(200)
	w.Write(append(responseBody, byte('\n')))
}

func (server *Server) GetAllPairsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	result, err := server.hashMap.GetAll()

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Error"))

		return
	}

	responseBody, _ := json.Marshal(result)

	w.WriteHeader(200)
	w.Write(append(responseBody, byte('\n')))
}

func (server *Server) GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))
}
