package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// NOTE: don't do this in real life
type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) add(w http.ResponseWriter, req *http.Request) {
	// some code here
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	// some code here
}

func (db database) fetch(w http.ResponseWriter, req *http.Request) {
	// some code here
}

func (db database) drop(w http.ResponseWriter, req *http.Request) {
	// some code here
}

func main() {
	db := database{
		"shoes": 50,
		"socks": 5,
	}

	// NOTE that these are all method values
	// (closing over the object "db")

	http.HandleFunc("/list", db.list)
	http.HandleFunc("/create", db.add)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.drop)
	http.HandleFunc("/read", db.fetch)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
