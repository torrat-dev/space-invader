package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Entity represents a game object's position on the board
type Entity struct {
	X int
	Y int
}

// GameState holds the entire state of the game
type GameState struct {
	Player         Entity
	Enemies        []Entity
	Bullets        []Entity
	EnemyBullets   []Entity
	Score          int
	GameOver       bool
	EnemyDirection int
	EnemySpeed     int
	EnemyFireCount int
}

var (
	state       GameState
	mu          sync.Mutex // Protects state during concurrent HTTP requests
	tmpl        *template.Template
	connections map[*websocket.Conn]bool
	connMu      sync.Mutex
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func init() {
	// Initialize game state with the player in the bottom middle
	// Assuming a 600px wide by 400px high game board
	state = GameState{
		Player:         Entity{X: 280, Y: 350},
		Bullets:        []Entity{},
		EnemyBullets:   []Entity{},
		EnemyDirection: 1,
		EnemySpeed:     5,
	}

	spawnEnemies()

	// Parse all HTML files in the templates directory
	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	connections = make(map[*websocket.Conn]bool)

	// Start the game loop at 60 FPS
	go func() {
		ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			updateGameState()
			broadcastState()
			mu.Unlock()
		}
	}()
}

func spawnEnemies() {
	state.Enemies = []Entity{}
	for row := 0; row < 5; row++ {
		for col := 0; col < 8; col++ {
			state.Enemies = append(state.Enemies, Entity{X: 50 + col*50, Y: 50 + row*30})
		}
	}
}

func updateGameState() {
	if !state.GameOver {
		// Update bullets
		newBullets := []Entity{}
		for _, b := range state.Bullets {
			b.Y -= 10
			if b.Y > 0 && b.Y < 400 {
				newBullets = append(newBullets, b)
			}
		}
		state.Bullets = newBullets

		// Update enemy bullets
		newEnemyBullets := []Entity{}
		for _, b := range state.EnemyBullets {
			b.Y += 8
			if b.Y < 400 {
				newEnemyBullets = append(newEnemyBullets, b)
			}
		}
		state.EnemyBullets = newEnemyBullets

		// Enemies fire bullets
		state.EnemyFireCount++
		if state.EnemyFireCount > 20 && len(state.Enemies) > 0 {
			state.EnemyFireCount = 0
			// Pick a random enemy to fire
			randomEnemy := state.Enemies[len(state.Enemies)-1]
			state.EnemyBullets = append(state.EnemyBullets, Entity{X: randomEnemy.X + 5, Y: randomEnemy.Y + 20})
		}

		// Update enemies
		if len(state.Enemies) > 0 {
			maxX := 0
			minX := 600
			for _, e := range state.Enemies {
				if e.X > maxX {
					maxX = e.X
				}
				if e.X < minX {
					minX = e.X
				}
			}
			if state.EnemyDirection == 1 && maxX >= 550 {
				state.EnemyDirection = -1
				for i := range state.Enemies {
					state.Enemies[i].Y += 20
				}
			} else if state.EnemyDirection == -1 && minX <= 0 {
				state.EnemyDirection = 1
				for i := range state.Enemies {
					state.Enemies[i].Y += 20
				}
			}
			for i := range state.Enemies {
				state.Enemies[i].X += state.EnemyDirection * state.EnemySpeed
			}
		}

		// Collision detection: bullets vs enemies
		for i := len(state.Bullets) - 1; i >= 0; i-- {
			b := state.Bullets[i]
			for j := len(state.Enemies) - 1; j >= 0; j-- {
				e := state.Enemies[j]
				if b.X-e.X <= 10 && e.X-b.X <= 10 && b.Y-e.Y <= 10 && e.Y-b.Y <= 10 {
					state.Bullets = append(state.Bullets[:i], state.Bullets[i+1:]...)
					state.Enemies = append(state.Enemies[:j], state.Enemies[j+1:]...)
					state.Score += 10
					break
				}
			}
		}

		// Collision detection: enemy bullets vs player
		for i := len(state.EnemyBullets) - 1; i >= 0; i-- {
			b := state.EnemyBullets[i]
			if b.X-state.Player.X <= 20 && state.Player.X-b.X <= 20 && b.Y-state.Player.Y <= 20 && state.Player.Y-b.Y <= 20 {
				state.GameOver = true
				break
			}
		}

		// Spawn new wave if no enemies left
		if len(state.Enemies) == 0 {
			spawnEnemies()
			state.EnemyDirection = 1
			state.EnemySpeed += 1
		}

		// Check game over
		for _, e := range state.Enemies {
			if e.Y >= 350 {
				state.GameOver = true
				break
			}
		}
	}
}

func resetState() {
	state = GameState{
		Player:         Entity{X: 280, Y: 350},
		Bullets:        []Entity{},
		EnemyBullets:   []Entity{},
		EnemyDirection: 1,
		EnemySpeed:     5,
		EnemyFireCount: 0,
	}
	spawnEnemies()
}

func broadcastState() {
	connMu.Lock()
	defer connMu.Unlock()
	mu.Lock()
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "game.html", state); err != nil {
		log.Println("Template execution error:", err)
		mu.Unlock()
		return
	}
	html := buf.String()
	mu.Unlock()
	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(html)); err != nil {
			conn.Close()
			delete(connections, conn)
		}
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/move-left", moveLeftHandler)
	http.HandleFunc("/move-right", moveRightHandler)
	http.HandleFunc("/shoot", shootHandler)
	http.HandleFunc("/score", scoreHandler)
	http.HandleFunc("/restart", restartHandler)

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

func moveLeftHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	state.Player.X -= 10
	if state.Player.X < 0 {
		state.Player.X = 0
	}
	w.WriteHeader(http.StatusOK)
}

func moveRightHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	state.Player.X += 10
	if state.Player.X > 570 { // Assuming spaceship width 30px, board 600px
		state.Player.X = 570
	}
	w.WriteHeader(http.StatusOK)
}

func shootHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	bullet := Entity{X: state.Player.X + 15, Y: state.Player.Y - 10}
	state.Bullets = append(state.Bullets, bullet)
	w.WriteHeader(http.StatusOK)
}

func scoreHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprintf(w, "%d", state.Score)
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	resetState()
	w.WriteHeader(http.StatusOK)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	connMu.Lock()
	connections[conn] = true
	connMu.Unlock()
	// Send initial state
	mu.Lock()
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "game.html", state); err != nil {
		log.Println("Template execution error:", err)
		mu.Unlock()
		return
	}
	html := buf.String()
	mu.Unlock()
	conn.WriteMessage(websocket.TextMessage, []byte(html))
}
