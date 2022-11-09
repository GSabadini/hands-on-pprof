package main

import (
	"database/sql"
	"encoding/json"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type response struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
}

func (r response) send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	json.NewEncoder(w).Encode(r.Message)
}

func SelectHandler(id int, db *sql.DB, trace trace.TracerProvider) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		result, err := Select(id, db, trace)
		if err != nil {
			response{http.StatusInternalServerError, err.Error()}.send(res)
		}

		response{http.StatusOK, result}.send(res)
	}
}

func InsertHandler(name string, db *sql.DB, trace trace.TracerProvider) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		result, err := Insert(name, db, trace)
		if err != nil {
			response{http.StatusInternalServerError, err.Error()}.send(res)
		}

		response{http.StatusCreated, result}.send(res)
	}
}
