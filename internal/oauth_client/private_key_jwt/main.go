// Package main. main is a sample application that demonstrates how to use the private_key_jwt method for OAuth 2.0 authentication with Casdoor.
// It is a simple web application that uses the private_key_jwt method to exchange an authorization code for an access token.
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/login", http.StatusFound)
	})
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)

	log.Println("Client running on", baseOauth2ClientAppUrl)
	log.Fatal(http.ListenAndServe(port, nil))
}
