package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const (
	scalderoneID = "aMfosX2ULwDaA3cc"
	scalderoneSecret = "0IsVLbFwm6DfsHJ2ytUaETmp9zz4I4pf"
	port = ":8080"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/auth", auth).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static"))).Methods("GET")
	fmt.Printf("Servre is running on localhost%s", port)
	panic(http.ListenAndServe(port, r))
}