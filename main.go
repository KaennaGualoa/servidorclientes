package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Customer struct {
	ID        string `json:"_id"`
	Nome      string `json:"nome"`
	Sobrenome string `json:"sobrenome"`
}

var customers []Customer

func main() {
	customers = make([]Customer, 0)

	r := mux.NewRouter()

	r.HandleFunc("/clients", GetClients).Methods("GET")
	r.HandleFunc("/clients", AddClient).Methods("POST")
	r.HandleFunc("/clients", UpdateClient).Methods("PUT")
	r.HandleFunc("/clients", DeleteClient).Methods("DELETE")

	http.Handle("/", r)

	fmt.Println("Servidor rodando na porta 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func GetClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func AddClient(w http.ResponseWriter, r *http.Request) {
	var client Customer
	_ = json.NewDecoder(r.Body).Decode(&client)
	if client.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "O campo '_id' não pode estar vazio.")
		return
	}
	if client.Nome == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "O campo 'nome' não pode estar vazio.")
		return
	}
	if client.Sobrenome == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "O campo 'sobrenome' não pode estar vazio.")
		return
	}
	for _, c := range customers {
		if c.ID == client.ID {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintln(w, "O usuário informado já está cadastrado.")
			return
		}
	}
	customers = append(customers, client)
	w.WriteHeader(http.StatusOK)
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {
	var updateData Customer
	_ = json.NewDecoder(r.Body).Decode(&updateData)

	for i, client := range customers {
		if client.ID == updateData.ID {
			c := customers[i]
			if updateData.Nome != "" {
				c.Nome = updateData.Nome
			}
			if updateData.Sobrenome != "" {
				c.Sobrenome = updateData.Sobrenome
			}
			customers[i] = c
			json.NewEncoder(w).Encode(customers[i])
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	var deleteData Customer
	_ = json.NewDecoder(r.Body).Decode(&deleteData)

	for i, client := range customers {
		if client.ID == deleteData.ID {
			customers = append(customers[:i], customers[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
