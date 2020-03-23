package main

import (
	"fmt"
	"log"
	"net/http"
)

// func serveHome(w http.ResponseWriter, r *http.Request) {
// 	log.Println(r.URL)

// 	if r.URL.Path != "/" {
// 		http.Error(w, "无", http.StatusNotFound)
// 		return
// 	}
// 	if r.Method != "GET" {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	http.ServeFile(w, r, "static/home.html")
// }

func main() {
	hub := newHub()
	//go hub.run()
	http.Handle("/", http.FileServer(http.Dir("static")))
	// TODO https://hack.chat/?your-channel
	//URL query param
	//all the params behind host is supposed to be split joint a particular channel
	http.HandleFunc("/chat-ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	fmt.Println("Listen on 6060")
	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
