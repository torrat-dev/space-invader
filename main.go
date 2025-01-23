package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// GameState holds the current state of the game
type GameState struct {
	Connections map[*websocket.Conn]string `json:"-"`
	Spaceship   Spaceship                  `json:"Spaceship"`
	Aliens      []Alien                    `json:"Aliens"`
	Projectiles []Projectile               `json:"Projectiles"`
}

// Spaceship represents the player's ship
type Spaceship struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// Alien represents an enemy alien
type Alien struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// Projectile represents a bullet fired by the spaceship
type Projectile struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// Hitbox represents a rectangular hitbox
type Hitbox struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// global variables
var upgrader = websocket.Upgrader{}
var gameState GameState
var alienWidth float64
var startX float64

// collidesWith checks if two hitboxes overlap
func (h Hitbox) collidesWith(other Hitbox) bool {
	return h.X+h.Width > other.X &&
		h.X < other.X+other.Width &&
		h.Y+h.Height > other.Y &&
		h.Y < other.Y+other.Height
}

func init() {
	alienWidth = 30.0     // Width of each alien
	alienSpacing := 100.0 // Desired spacing between aliens
	numAliensPerRow := 4  // Number of aliens per row
	gameAreaWidth := 600.0

	totalAlienWidth := float64(numAliensPerRow) * alienWidth
	totalSpacingWidth := float64(numAliensPerRow-1) * alienSpacing
	availableWidth := gameAreaWidth - totalAlienWidth - totalSpacingWidth
	startX = availableWidth / 2 // Start position to center the aliens

	// Initialize the game state
	gameState = GameState{
		Connections: make(map[*websocket.Conn]string),
		Spaceship:   Spaceship{X: 275, Y: 360, Width: 50, Height: 20},
		Aliens:      make([]Alien, 0),
		Projectiles: make([]Projectile, 0),
	}

	// Initialize aliens
	for i := 0; i < 4; i++ {
		for j := 0; j < 5; j++ {
			x := startX + float64(i)*(alienWidth+alienSpacing)
			gameState.Aliens = append(gameState.Aliens, Alien{X: x, Y: float64(j*30 + 30), Width: alienWidth, Height: 20})
		}
	}
}

func moveSpaceship(direction int) {
	// Update spaceship position
	gameState.Spaceship.X += float64(direction * 10) // Adjust speed as needed

	// Keep spaceship within bounds
	if gameState.Spaceship.X < 0 {
		gameState.Spaceship.X = 0
	} else if gameState.Spaceship.X > 550 { // Assuming game area width is 600 and spaceship width is 50
		gameState.Spaceship.X = 550
	}
}

func gameLoop() {
	for {
		// Update game state (move aliens, projectiles, etc.)
		moveAliens()
		moveProjectiles()
		checkCollisions()

		for c, connectionId := range gameState.Connections {

			if err := c.WriteMessage(websocket.TextMessage, generateGameStateJSON()); err != nil {
				log.Printf("\n Game loop connection error: %s connectionId: %s", err, connectionId)
				delete(gameState.Connections, c)
			}
		}

		time.Sleep(time.Millisecond * 16) // Approximately 60 updates per second
	}
}

func moveAliens() {
	// Adjust these values to control alien movement
	alienSpeed := 0.25  // Speed of alien movement
	alienDirection := 1 // 1 for right, -1 for left

	// Move each alien
	for i := range gameState.Aliens {
		gameState.Aliens[i].X += float64(alienDirection) * alienSpeed

		// Check if aliens hit the edge of the game area
		if gameState.Aliens[i].X+gameState.Aliens[i].Width > 600 || gameState.Aliens[i].X < 0 {
			// Reverse direction and move down
			alienDirection *= -1
			for j := range gameState.Aliens {
				gameState.Aliens[j].Y += 10 // Move aliens down
			}
			gameState.Aliens[i].X = startX
			gameState.Aliens[i].Width = alienWidth
			break // Exit the loop after reversing direction
		}
	}
}

func moveProjectiles() {
	// Implement projectile movement logic here
	// ...
}

func checkCollisions() {
	// Implement collision detection logic here
	// ...
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {

		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		// Handle incoming messages (e.g., player input)
		var input struct {
			Direction int `json:"direction"`
		}

		if err = json.Unmarshal(message, &input); err != nil {
			log.Println("unmarshal:", err)
			break
		}

		moveSpaceship(input.Direction)

		connectionId := uuid.New().String()

		gameState.Connections[c] = connectionId // Send updated game state to the client

		err = c.WriteMessage(mt, generateGameStateJSON())
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func generateGameStateJSON() []byte {
	// Generate JSON representation of the game state
	gameStateJSON, err := json.Marshal(gameState)
	if err != nil {
		log.Println("marshal:", err)
		return []byte{}
	}
	return gameStateJSON
}

func main() {
	// Start the game loop in a separate goroutine
	go gameLoop()

	http.HandleFunc("/ws", socketHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
