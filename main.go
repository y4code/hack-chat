// package main

// import (
// 	"log"
// 	"net/http"
// )

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
// 	http.ServeFile(w, r, "./static/home.html")
// }

package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8080", nil)
}
