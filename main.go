package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	//if r.URL.Path != "/" {
	//	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	//	//http.Error(w, "Not Found", http.StatusNotFound)
	//	return
	//}
	//if r.Method != "GET" {
	//	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	//	return
	//}
	//http.ServeFile(w, r, "home.html")
}


func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
