// main
package main

import (
	"fmt"
	"log"
	"net/http"
)

func midHandler1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Top Handler Wrapper - Start")
		next.ServeHTTP(w, r)
		log.Println("Top Handler Wrapper - End")
	})
}
func midHandler2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("2nd Handler Wrapped By Top - Started")
		if r.URL.Path == "/message" {
			if r.URL.Query().Get("password") == "password" {
				log.Println("Access Authorized...")
				next.ServeHTTP(w, r)
			} else {
				log.Println("Failed to authorize to the system")
				return
			}
		} else {
			next.ServeHTTP(w, r)
		}
		log.Println("2nd Handler Wrapped By Top - Ended")
	})
}
func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing index Handler")
	fmt.Fprintf(w, "Welcome")
}
func message(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing message Handler")
	fmt.Fprintf(w, "HTTP Middleware is awesome")
}
func iconHandler(w http.ResponseWriter, r *http.Request) {
}
func main() {
	http.HandleFunc("/favicon.ico", iconHandler)
	http.Handle("/", midHandler1(midHandler2(http.HandlerFunc(index))))
	http.Handle("/message", midHandler1(midHandler2(http.HandlerFunc(message))))
	server := &http.Server{
		Addr: ":8080",
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
