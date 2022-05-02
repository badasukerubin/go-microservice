package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello World")
		d, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Oops", http.StatusBadRequest)
			return
		}
		log.Printf("Data %s", d)
	})

	http.ListenAndServe(":9090", nil)
}
