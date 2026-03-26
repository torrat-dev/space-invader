package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Entity represents a game object's position on the board
type Entity struct {
	X int
	Y int
}

// GameState holds the entire state of the game
type GameState struct {
	Player   Entity
	Enemies  []Entity
	Bullets  []Entity
	Score    int
	GameOver bool
}

var (
	state GameState
	mu    sync.Mutex // Protects state during concurrent HTTP requests
	tmpl  *template.Template
)

func init() {
	// Initialize game state with the player in the bottom middle
	// Assuming a 600px wide by 400px high game board
	state = GameState{
		Player: Entity{X: 280, Y: 20},
	}

	// Parse all HTML files in the templates directory
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/game", gameHandler)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Pass the current state to the game board template
	if err := tmpl.ExecuteTemplate(w, "game.html", state); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
