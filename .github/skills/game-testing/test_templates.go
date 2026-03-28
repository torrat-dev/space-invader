// Test templates for space-invader game functionality
// Use these as starting points for writing tests in main_test.go

import (
	"sync"
	"testing"
)

// Template for testing player movement
func TestPlayerMovement(t *testing.T) {
	mu := &sync.Mutex{}
	state := GameState{Player: Entity{X: 280, Y: 350}}

	// Test move left
	mu.Lock()
	state.Player.X -= 10
	if state.Player.X < 0 {
		state.Player.X = 0
	}
	mu.Unlock()
	// Assertions: check position

	// Test move right
	mu.Lock()
	state.Player.X += 10
	if state.Player.X > 570 {
		state.Player.X = 570
	}
	mu.Unlock()
	// Assertions
}

// Template for testing bullet shooting
func TestBulletShooting(t *testing.T) {
	mu := &sync.Mutex{}
	state := GameState{Player: Entity{X: 280, Y: 350}, Bullets: []Entity{}}

	mu.Lock()
	bullet := Entity{X: state.Player.X + 15, Y: state.Player.Y - 10}
	state.Bullets = append(state.Bullets, bullet)
	mu.Unlock()
	// Assertions: check bullet added
}

// Template for testing enemy movement
func TestEnemyMovement(t *testing.T) {
	mu := &sync.Mutex{}
	state := GameState{Enemies: []Entity{{X: 50, Y: 50}}, EnemyDirection: 1, EnemySpeed: 5}

	mu.Lock()
	for i := range state.Enemies {
		state.Enemies[i].X += state.EnemyDirection * state.EnemySpeed
	}
	mu.Unlock()
	// Assertions: check new positions
}

// Template for testing collision detection
func TestCollisionDetection(t *testing.T) {
	mu := &sync.Mutex{}
	state := GameState{
		Bullets: []Entity{{X: 50, Y: 50}},
		Enemies: []Entity{{X: 50, Y: 50}},
		Score:   0,
	}

	mu.Lock()
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
	mu.Unlock()
	// Assertions: check score increased, entities removed
}

// Template for testing game over conditions
func TestGameOver(t *testing.T) {
	mu := &sync.Mutex{}
	state := GameState{Enemies: []Entity{{X: 50, Y: 360}}, GameOver: false}

	mu.Lock()
	for _, e := range state.Enemies {
		if e.Y >= 350 {
			state.GameOver = true
			break
		}
	}
	mu.Unlock()
	// Assertions: check GameOver is true
}

// Add more templates as needed for WebSocket, HTTP handlers, etc.