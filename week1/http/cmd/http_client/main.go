package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"
)

const (
	baseUrl       = "http://localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/%d"
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

func createNote() (Note, error) {
	note := NoteInfo{
		Title:    gofakeit.BeerName(),
		Context:  gofakeit.IPv4Address(),
		Author:   gofakeit.Name(),
		IsPublic: gofakeit.Bool(),
	}
	data, err := json.Marshal(note)
	if err != nil {
		return Note{}, err
	}
	resp, err := http.Post(baseUrl+createPostfix, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Note{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return Note{}, err
	}

	var createNote Note

	if err := json.NewDecoder(resp.Body).Decode(&createNote); err != nil {
		return Note{}, err
	}

	return createNote, nil
}

func getNote(id int64) (Note, error) {
	url := fmt.Sprintf(baseUrl+getPostfix, id)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return Note{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Note{}, fmt.Errorf("failed getting note by note_id:%d, status code:%d", id, resp.StatusCode)
	}

	var respNote Note
	if err := json.NewDecoder(resp.Body).Decode(&respNote); err != nil {
		return Note{}, err
	}
	return respNote, nil
}

func main() {
	note, err := createNote()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(color.RedString("Note created:\n"), color.GreenString("%+v", note))
	note, err = getNote(note.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(color.RedString("Note getting:\n"), color.GreenString("%+v", note))
}
