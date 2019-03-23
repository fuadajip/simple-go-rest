package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	// Bird ...
	Bird struct {
		Species     string `json:"species"`
		Description string `json:"description"`
	}
)

var birds []Bird

func getBirdHandler(w http.ResponseWriter, r *http.Request) {
	// convert "birds" variable to json
	birdListBytes, err := json.Marshal(birds)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(birdListBytes)
}

func createBirdHandler(w http.ResponseWriter, r *http.Request) {
	// create new instance of bird
	bird := Bird{}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error:  %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bird.Species = r.Form.Get("species")
	bird.Description = r.Form.Get("description")

	// appen bird data to existing list as temporary db
	birds = append(birds, bird)

	http.Redirect(w, r, "/assets/", http.StatusFound)
}
