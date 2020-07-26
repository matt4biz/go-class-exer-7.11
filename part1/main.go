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
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if _, ok := db[item]; ok {
		w.WriteHeader(http.StatusBadRequest) // 404

		fmt.Fprintf(w, "duplicate item: %q\n", item)
		return
	}

	if f64, err := strconv.ParseFloat(price, 32); err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400

		fmt.Fprintf(w, "invalid price: %q\n", price)
	} else {
		db[item] = dollars(f64)

		fmt.Fprintf(w, "added %s with price %s\n", item, dollars(f64))
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusNotFound) // 404

		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	if f64, err := strconv.ParseFloat(price, 32); err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400

		fmt.Fprintf(w, "invalid price: %q\n", price)
	} else {
		db[item] = dollars(f64)

		fmt.Fprintf(w, "new price %s for %s\n", dollars(f64), item)
	}
}

func (db database) fetch(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusNotFound) // 404

		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	fmt.Fprintf(w, "item %s has price %s\n", item, db[item])
}

func (db database) drop(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusNotFound) // 404

		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	delete(db, item)

	fmt.Fprintf(w, "dropped %s\n", item)
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
