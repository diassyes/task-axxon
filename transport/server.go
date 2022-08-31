package transport

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"task-axxon/config"
	"task-axxon/model"
	"task-axxon/store"
	"task-axxon/view"
)

func InitServer(conf config.Config) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/", ProxyHTTP).Methods(http.MethodPost)
	router.HandleFunc("/", GetByUUID).Methods(http.MethodGet)
	return &http.Server{
		Addr:    conf.Address,
		Handler: router,
	}
}

func GetByUUID(w http.ResponseWriter, req *http.Request) {
	data := store.DefaultStore.Get(req.URL.Query().Get("uuid"))
	if data == nil {
		http.Error(w, "uuid not found", http.StatusNotFound)
		return
	}
	_, err := w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ProxyHTTP(w http.ResponseWriter, reqFromClient *http.Request) {
	var (
		viewReq  view.Request
		viewResp view.Response
	)
	err := json.NewDecoder(reqFromClient.Body).Decode(&viewReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(viewReq.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reqToOut, err := http.NewRequest(viewReq.Method, viewReq.URL, b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if viewReq.Headers != nil {
		reqToOut.Header = viewReq.Headers
	}
	respFromOut, err := http.DefaultClient.Do(reqToOut)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer respFromOut.Body.Close()

	id := uuid.New().String()
	viewResp = view.NewResponse(respFromOut, id)
	modelReq := model.NewRequest(&viewReq, &viewResp)
	dataToSave, err := json.Marshal(modelReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	store.DefaultStore.Set(id, dataToSave)
	dataToOut, err := json.Marshal(viewResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	_, err = w.Write(dataToOut)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}
