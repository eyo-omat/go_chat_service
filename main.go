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

func getRandomName() string {
	adjs := []string{"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark", "summer", "icy", "delicate", "quiet", "white", "cool", "spring", "winter", "patient", "twilight", "dawn", "crimson", "wispy", "weathered", "blue", "billowing", "broken", "cold", "damp", "falling", "frosty", "green", "long", "late", "lingering", "bold", "little", "morning", "muddy", "old", "red", "rough", "still", "small", "sparkling", "throbbing", "shy", "wandering", "withered", "wild", "black", "young", "holy", "solitary", "fragrant", "aged", "snowy", "proud", "floral", "restless", "divine", "polished", "ancient", "purple", "lively", "nameless"}
	nouns := []string{"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning", "snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter", "forest", "hill", "cloud", "meadow", "sun", "glade", "bird", "brook", "butterfly", "bush", "dew", "dust", "field", "fire", "flower", "firefly", "feather", "grass", "haze", "mountain", "night", "pond", "darkness", "snowflake", "silence", "sound", "sky", "shape", "surf", "thunder", "violet", "water", "wildflower", "wave", "water", "resonance", "sun", "wood", "dream", "cherry", "tree", "fog", "frost", "voice", "paper", "frog", "smoke", "star"}
	return adjs[rand.Intn(len(adjs))] + "_" + nouns[rand.Intn(len(nouns))]
}

func getRandomColor() string {
	return "#" + strconv.FormatInt(rand.Int63n(0xFFFFFF), 16)
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