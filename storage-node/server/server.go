package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"auxdev.at/storage-server/storage"
	"github.com/gorilla/mux"
)

type MessageResponse struct {
	Message string `json:"message"`
}

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

		w.WriteHeader(400)
		response, _ := json.Marshal(&MessageResponse{Message: "Invalid request body."})
		w.Write(response)

		return
	}

	if err := server.hashMap.Put(pair); err != nil {
		server.logger.Println(err)

		w.WriteHeader(500)
		response, _ := json.Marshal(&MessageResponse{Message: "Internal server error while trying to store pair."})
		w.Write(response)

		return
	}

	server.logger.Printf("Stored %s", pair)

	defer r.Body.Close()

	w.WriteHeader(200)
	response, _ := json.Marshal(&MessageResponse{Message: "OK"})
	w.Write(response)
}

func (server *Server) GetPairHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	params := mux.Vars(r)
	key := params["key"]

	result, err := server.hashMap.Get(key)

	if err != nil {
		server.logger.Println(err)

		w.WriteHeader(404)
		response, _ := json.Marshal(&MessageResponse{Message: "Not found."})
		w.Write(response)

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
		w.WriteHeader(500)
		response, _ := json.Marshal(&MessageResponse{Message: "Internal error while trying to load all keys"})
		w.Write(response)

		return
	}

	responseBody, _ := json.Marshal(result)

	w.WriteHeader(200)
	w.Write(append(responseBody, byte('\n')))
}

func (server *Server) GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	response, _ := json.Marshal(&MessageResponse{Message: "OK"})
	w.Write(response)
}
