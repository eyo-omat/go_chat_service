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

type customClaims struct {
	jwt.StandardClaims
	Client string `json:"client"`
	Channel string `json:"channel"`
	Data userData `json:"data"`
	Permissions map[string] permissionClaims `json:"permissions"`
}

type permissionClaims struct {
	Publish	bool	`json:"publish"`
	Subscribe	bool	`json:"subscribe"`
}

type userData struct {
	Color string `json:"color"`
	Name string `json:"name"`
}

func auth (w http.ResponseWriter, r *http.Request)  {
	clientID := r.FormValue("clientID")
	if	clientID == "" {
		http.Error(w, "No client ID defined", http.StatusUnprocessableEntity)
		return
	}

	// Public Room
	publicRoomRegex := "^observable-room$"
	// Private Room of the request user
	userPrivateRoomRegex := fmt.Sprintf("^private-room-%s$", clientID)
	// Private Rooms of other users except request user
	otherUsersPrivateRoomRegex := fmt.Sprintf("^private-room-(?!%s$).+$", clientID)

	claims := customClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 3).Unix(),
		},
		Client: clientID,
		Channel: scalderoneID,
		Data: userData{
			Name: getRandomName(),
			Color: getRandomColor(),
		},
		Permissions: map[string]permissionClaims{
			publicRoomRegex: permissionClaims{
				Publish: true,
				Subscribe: true,
			},
			userPrivateRoomRegex: permissionClaims{
				Publish: false,
				Subscribe: true,
			},
			otherUsersPrivateRoomRegex: permissionClaims{
				Publish: true,
				Subscribe: false,
			},
		},
	}

	// Generate a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(scalderoneSecret))
	if err != nil {
		http.Error(w, "Unable to sign token", http.StatusUnprocessableEntity)
		return
	}

	//Send to the user the token
	w.Write([]byte(tokenString))
}