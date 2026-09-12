package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Note struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var notes = make([]Note, 0)
var nextID = 1

func notesHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	case http.MethodPost:
		var note Note
		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		note.ID = nextID
		nextID++
		notes = append(notes, note)

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func noteByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idstr := strings.TrimPrefix(r.URL.Path, "/notes/")
		id, err := strconv.Atoi(idstr)

		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		for _, note := range notes {
			if note.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(note)
				return
			}
			http.Error(w, "note not found", http.StatusNotFound)
		}
	case http.MethodPut:
		idstr := strings.TrimPrefix(r.URL.Path, "/notes/")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		var updatedNote Note

		err = json.NewDecoder(r.Body).Decode(&updatedNote)

		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		for i := range notes {

			if notes[i].ID == id {
				notes[i].Title = updatedNote.Title
				notes[i].Description = updatedNote.Description

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(notes[i])
				return
			}
			http.Error(w, "note not found", http.StatusNotFound)
		}
	case http.MethodDelete:
		idstr := strings.TrimPrefix(r.URL.Path, "/notes/")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		for i := range notes {
			if notes[i].ID == id {
				notes = append(notes[:i], notes[i+1:]...)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		http.Error(w, "note not found", http.StatusNotFound)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/notes", notesHandler)
	http.HandleFunc("/notes/", noteByIDHandler)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
