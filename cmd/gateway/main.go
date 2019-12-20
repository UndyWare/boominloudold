package gateway

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Data struct with just an id and a body tag
type Data struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

var dataset = []Data{
	Data{"1", "bruh.mp3"},
	Data{"2", "theWATCH.mp3"},
	Data{"3", "haroldthealien.mp3"},
	Data{"4", "schmeat.mp3"},
	Data{"5", "distortedbass.mp3"},
	Data{"6", "yeet.mp3"},
	Data{"7", "whoadance.mp3"},
	Data{"8", "goodolerub.mp3"}}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataset)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parameters := mux.Vars(r)

	for _, item := range dataset {
		if item.ID == parameters["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	// I believe this should "return" an empty Data struct since this will only be reached
	// in the case of the specified data not being found
	json.NewEncoder(w).Encode(&Data{})
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data Data
	_ = json.NewDecoder(r.Body).Decode(&data)
	data.ID = strconv.Itoa(rand.Intn(100000))
	dataset = append(dataset, data)
	json.NewEncoder(w).Encode(&data)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(r)
	for index, item := range dataset {
		if item.ID == parameters["id"] {
			// Remove the old version of this item so we can update it
			dataset = append(dataset[:index], dataset[index+1:]...)

			var data Data
			_ = json.NewDecoder(r.Body).Decode(&data)
			data.ID = parameters["id"]
			dataset = append(dataset, data)
			break
		}
	}
	json.NewEncoder(w).Encode(dataset)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(r)
	for index, item := range dataset {
		if item.ID == parameters["id"] {
			// We will delete by just creating a new slice containing everything but item
			dataset = append(dataset[:index], dataset[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(dataset)
}

// Gateway init
func Gateway() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	router := mux.NewRouter()

	router.HandleFunc("/sounds", getPosts).Methods("GET")
	router.HandleFunc("/sounds", createPost).Methods("POST")
	router.HandleFunc("/sounds/{id}", getPost).Methods("GET")
	router.HandleFunc("/sounds/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/sounds/{id}", deletePost).Methods("DELETE")

	// Now to serve static files:
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Successfully started gateway at port:", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
