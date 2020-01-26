package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home landing for customers.")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	router.HandleFunc("/customer", createCustomer).Methods("POST")
	router.HandleFunc("/customer/{id}", getOneCustomer).Methods("GET")
	router.HandleFunc("/customer/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customers", getAllCustomers).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type customer struct {
	ID          string `json:"ID"`
	Name       	string `json:"Name"`
	Email 		string `json:"Email"`
	Phone 		string `json:"Phone"`
}

// Use slice as a tempDB
type allCustomers []customer

var customers = allCustomers{
	{
		ID:         "1",
		Name:       "Bill Testperson",
		Email: 		"bill@example.com",
		Phone: 		"406-555-5555",
	},
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var newCustomer customer
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please enter data needed: Name, Email, Phone")
	}
	
	json.Unmarshal(reqBody, &newCustomer)
	customers = append(customers, newCustomer)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newCustomer)
}

func getOneCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]

	for _, singleCustomer := range customers {
		if singleCustomer.ID == customerID {
			json.NewEncoder(w).Encode(singleCustomer)
		}
	}
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	customerId := mux.Vars(r)["id"]
	var updatedCustomer customer

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please enter data needed: Name, Email, Phone")
	}
	json.Unmarshal(reqBody, &updatedCustomer)

	for i, singleCustomer := range customers {
		if singleCustomer.ID == customerId {
			singleCustomer.Name = updatedCustomer.Name
			singleCustomer.Email = updatedCustomer.Email
			singleCustomer.Phone = updatedCustomer.Phone
			customers = append(customers[:i], singleCustomer)
			json.NewEncoder(w).Encode(singleCustomer)
		}
	}
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(customers)
}
