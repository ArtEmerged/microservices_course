package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
)

const (
	baseUrl       = "localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/{id}"
)

type NoteInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"is_public"`
}

type Note struct {
	ID        int64     `json:"id"`
	Info      NoteInfo  `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SyncMap struct {
	incrementID int64
	elems       map[int64]*Note //id:Note
	m           sync.RWMutex
}

var notes = &SyncMap{
	elems: make(map[int64]*Note),
	// мьютекс можно не иницальизировать, он по дефолту уже есть
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var (
		info = &NoteInfo{}
		err  error
	)
	if err = json.NewDecoder(r.Body).Decode(info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	now := time.Now()

	note := &Note{
		ID:        notes.incrementID,
		Info:      *info,
		CreatedAt: now,
		UpdatedAt: now,
	}
	notes.incrementID++

	w.Header().Set("Content-Type", "applicatin/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "failed to encoded note data", http.StatusInternalServerError)
		return
	}

	notes.m.Lock()
	defer notes.m.Unlock()
	notes.elems[note.ID] = note

}

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(noteID, 10, 64)
	if err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	notes.m.RLock()
	defer notes.m.RUnlock()
	note, ok := notes.elems[id]
	if !ok {
		http.Error(w, "note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "failed to encoded note data", http.StatusInternalServerError)
	}
}

func main() {
	r := chi.NewRouter()
	r.Get(getPostfix, getNoteHandler)
	r.Post(createPostfix, createNoteHandler)
	fmt.Println("Listen", baseUrl)
	if err := http.ListenAndServe(baseUrl, r); err != nil {
		log.Fatal(err)
	}

}
