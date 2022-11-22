package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// The struct for Pokemon data
type Pokemon struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	WillEvolve bool   `json:"willEvolve"`
}

// some data to start - TODO: implement DB
//
//	var Pokedex = []Pokemon{
//		{ID: "001", Name: "Grookey", Type: "Grass", WillEvolve: true},
//		{ID: "002", Name: "Thwackey", Type: "Grass", WillEvolve: true},
//		{ID: "003", Name: "Rillaboom", Type: "Grass", WillEvolve: false},
//	}
var Pokedex []Pokemon

func main() {
	// create a new router
	router := mux.NewRouter()
	populatePokedex()
	// specify endpoints, handler functions and HTTP method
	router.HandleFunc("/pokedex", getPokedex).Methods("GET")
	router.HandleFunc("/pokedex", addPokemon).Methods("POST")
	router.HandleFunc("/pokedex/{id}", getPokemonByID).Methods("GET")
	fmt.Println("Pokedex app is up and running")

	// start and listen to requests
	http.ListenAndServe(":8080", router)
}

// getPokedex responds with the list of all pokemon as JSON
func getPokedex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Show the entire Pokedex")
	// return the Pokedex slice as JSON
	json.NewEncoder(w).Encode(Pokedex)
}

// addPokemon adds a pokemon received by a JSON request body
func addPokemon(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// define a pokemon variable with Pokemon as type
	var newPokemon Pokemon

	json.NewDecoder(r.Body).Decode(&newPokemon)

	for _, pokemon := range Pokedex {
		if pokemon.ID == newPokemon.ID {
			fmt.Println("This Pokemon is already in the Pokedex")
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	// append it to the Pokedex slice
	Pokedex = append(Pokedex, newPokemon)

	bytePokedex, err := json.Marshal(Pokedex)
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("containerData/pokemon.json", bytePokedex, 0644)

	fmt.Printf("%v is just added to the Pokedex\n", newPokemon.Name)
	json.NewEncoder(w).Encode(newPokemon)

}

// getPokemonByID show the data for the only requested pokemon if
// it's available in the pokedex
func getPokemonByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// retrieve the pokemon id from the route using mux.Vars()
	id := vars["id"]

	// iterate through the slice and return only the requested pokemon
	for _, pokemon := range Pokedex {
		if pokemon.ID == id {
			json.NewEncoder(w).Encode(pokemon)
		}
	}
}

func populatePokedex() {
	// get the file from the volume
	file, err := os.Open("containerData/pokemon.json")
	if err != nil {
		fmt.Println(err)
	}
	// defer the command to close the file
	defer file.Close()
	// read the file
	bytePokedex, err2 := io.ReadAll(file)
	if err2 != nil {
		fmt.Println(err2)
	}
	// unmarshal the byte slice and store in the Pokedex var
	json.Unmarshal(bytePokedex, &Pokedex)
}
