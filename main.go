package main

import (
	"flag"
	"log"
	"net/http"

)

var addr = flag.String("addr", ":8080", "http service	address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "client/index.html")
}

func main() {
	flag.Parse()
	hub := NewHub() 
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}



