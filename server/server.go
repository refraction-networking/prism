package main

import (
	"log"
	//"github.com/refraction-networking/Metis/endpoint"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"errors"
	"fmt"
	"strings"
)

type Website struct {
	Domain string `json:"domain,omitempty"`
	/*
	//Accuracy - number of times we've tested this domain and run into a problem indicative of censorship.
	//If we ever test this site and get through, we remove it from the blocked list.
	acc float `json:"???"
	 */
}

var blockedList []Website

func contains(slice []Website, s string) bool {
	for _, e := range slice {
		if strings.Contains(s, e.Domain) { return true}
	}
	return false
}

func getBlocked(writer http.ResponseWriter, req *http.Request) {
	//Creates a json encoder that writes to writer
	json.NewEncoder(writer).Encode(blockedList)
}

func addBlocked(writer http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		log.Fatal(errors.New("POST request from Metis client was empty"))
		return
	}
	dec := json.NewDecoder(req.Body)
	// read open bracket
	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	// while the array contains values
	for dec.More() {
		var d string
		// decode an array value (Message)
		err := dec.Decode(&d)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v \n", d)
		if !contains(blockedList, d) {
			blockedList = append(blockedList, Website{Domain: d})
		}
	}

	// read closing bracket
	t, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	//endpt := new(endpoint.Endpoint)
	log.Println("Starting Metis server...")
	router := mux.NewRouter()
	router.HandleFunc("/blocked", getBlocked).Methods("GET")
	router.HandleFunc("/blocked/add", addBlocked).Methods("POST")
	blockedList = append(blockedList, Website{Domain: "facebook.com"})
	blockedList = append(blockedList, Website{Domain: "google.com"})
	blockedList = append(blockedList, Website{Domain: "bettermotherfuckingwebsite.com"})
	log.Fatal(http.ListenAndServe(":9099", router))
}
